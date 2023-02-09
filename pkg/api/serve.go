package api

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.5.6 -package api -generate "types,client,chi-server,spec" -templates tmpl -o lakefs.gen.go ../../api/swagger.yml

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	lakefsv1 "github.com/treeverse/lakefs/gen/proto/go/lakefs/v1"
	"github.com/treeverse/lakefs/pkg/api/params"
	"github.com/treeverse/lakefs/pkg/auth"
	"github.com/treeverse/lakefs/pkg/auth/email"
	authoidc "github.com/treeverse/lakefs/pkg/auth/oidc"
	"github.com/treeverse/lakefs/pkg/block"
	"github.com/treeverse/lakefs/pkg/catalog"
	"github.com/treeverse/lakefs/pkg/cloud"
	"github.com/treeverse/lakefs/pkg/config"
	"github.com/treeverse/lakefs/pkg/httputil"
	"github.com/treeverse/lakefs/pkg/logging"
	"github.com/treeverse/lakefs/pkg/stats"
	"github.com/treeverse/lakefs/pkg/templater"
	"github.com/treeverse/lakefs/pkg/upload"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

const (
	RequestIDHeaderName = "X-Request-ID"
	LoggerServiceName   = "rest_api"
	BaseURL             = "/api/v1"

	extensionValidationExcludeBody = "x-validation-exclude-body"
)

func Serve(cfg *config.Config, catalog catalog.Interface, middlewareAuthenticator auth.Authenticator, controllerAuthenticator auth.Authenticator, authService auth.Service, blockAdapter block.Adapter, metadataManager auth.MetadataManager, migrator Migrator, collector stats.Collector, cloudMetadataProvider cloud.MetadataProvider, actions actionsHandler, auditChecker AuditChecker, logger logging.Logger, emailer *email.Emailer, templater templater.Service, gatewayDomains []string, snippets []params.CodeSnippet, oidcProvider *oidc.Provider, oauthConfig *oauth2.Config, pathProvider upload.PathProvider) http.Handler {
	logger.Info("initialize OpenAPI server")
	swagger, err := GetSwagger()
	if err != nil {
		panic(err)
	}

	sessionStore := sessions.NewCookieStore(authService.SecretStore().SharedSecret())
	r := chi.NewRouter()
	oidcConfig := cfg.Auth.OIDC
	apiRouter := r.With(
		OapiRequestValidatorWithOptions(swagger, &openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		}),
		httputil.LoggingMiddleware(
			RequestIDHeaderName,
			logging.Fields{logging.ServiceNameFieldKey: LoggerServiceName},
			cfg.Logging.AuditLogLevel,
			cfg.Logging.TraceRequestHeaders),
		AuthMiddleware(logger, swagger, middlewareAuthenticator, authService, sessionStore, &oidcConfig),
		MetricsMiddleware(swagger),
	)
	oidcAuthenticator := authoidc.NewAuthenticator(oauthConfig, oidcProvider)
	controller := NewController(
		cfg,
		catalog,
		controllerAuthenticator,
		authService,
		blockAdapter,
		metadataManager,
		migrator,
		collector,
		cloudMetadataProvider,
		actions,
		auditChecker,
		logger,
		emailer,
		templater,
		oidcAuthenticator,
		sessionStore,
		pathProvider,
	)
	HandlerFromMuxWithBaseURL(controller, apiRouter, BaseURL)

	r.Mount("/_health", httputil.ServeHealth())
	r.Mount("/metrics", promhttp.Handler())
	r.Mount("/_pprof/", httputil.ServePPROF("/_pprof/"))
	r.Mount("/swagger.json", http.HandlerFunc(swaggerSpecHandler))
	r.Mount(BaseURL, http.HandlerFunc(InvalidAPIEndpointHandler))
	if cfg.Auth.OIDC.Enabled {
		r.Mount("/oidc/login", NewOIDCLoginPageHandler(oidcConfig, sessionStore, oauthConfig, logger))
	}
	r.Mount("/logout", NewLogoutHandler(sessionStore, logger, cfg.Auth.LogoutRedirectURL))

	// Configuration flag to control if the embedded UI is served
	// or not and assign the correct handler for each case.
	var rootHandler http.Handler
	if cfg.UI.Enabled {
		// Handler which serves the embedded UI
		// as well as handles erroneous S3 gateway requests
		// and returns a compatible response
		rootHandler = NewUIHandler(gatewayDomains, snippets)
	} else {
		// Handler which only handles erroneous S3 gateway requests
		// and returns a compatible response
		rootHandler = NewS3GatewayEndpointErrorHandler(gatewayDomains)
	}

	service := &Service{
		Catalog:      controller.Catalog,
		BlockAdapter: controller.BlockAdapter,
	}

	grpcServer := grpc.NewServer()
	lakefsv1.RegisterLakeFSServiceServer(grpcServer, service)
	r.Mount("/", rootHandler)

	listenOn := ":8080"
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.ProtoMajor == 2 && strings.HasPrefix(request.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(writer, request)
		} else {
			r.ServeHTTP(writer, request)
		}
	})
}

func swaggerSpecHandler(w http.ResponseWriter, _ *http.Request) {
	reader, err := GetSwaggerSpecReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = io.Copy(w, reader)
}

// OapiRequestValidatorWithOptions Creates middleware to validate request by swagger spec.
// This middleware is good for net/http either since go-chi is 100% compatible with net/http.
// The original implementation can be found at https://github.com/deepmap/oapi-codegen/blob/master/pkg/chi-middleware/oapi_validate.go
// Used our own implementation in order to:
//  1. Use the latest version kin-openapi (can switch back when oapi-codegen will be updated)
//  2. For file upload wanted to skip body validation for two reasons:
//     a. didn't find a way for the validator to accept any file content type
//     b. didn't want the validator to read the complete request body for the specific request
func OapiRequestValidatorWithOptions(swagger *openapi3.Swagger, options *openapi3filter.Options) func(http.Handler) http.Handler {
	router, err := legacy.NewRouter(swagger)
	if err != nil {
		panic(err)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// find route
			route, m, err := router.FindRoute(r)
			if err != nil {
				// We failed to find a matching route for the request.
				writeError(w, r, http.StatusBadRequest, err.Error())
				return
			}

			// include operation id from route in the context for logging
			r = r.WithContext(logging.AddFields(r.Context(), logging.Fields{"operation_id": route.Operation.OperationID}))

			// validate request
			statusCode, err := validateRequest(r, route, m, options)
			if err != nil {
				writeError(w, r, statusCode, err.Error())
				return
			}
			// serve
			next.ServeHTTP(w, r)
		})
	}
}

func validateRequest(r *http.Request, route *routers.Route, pathParams map[string]string, options *openapi3filter.Options) (int, error) {
	// Extension - validation exclude body
	if _, ok := route.Operation.Extensions[extensionValidationExcludeBody]; ok {
		o := *options
		o.ExcludeRequestBody = true
		options = &o
	}

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
		Options:    options,
	}
	if err := openapi3filter.ValidateRequest(r.Context(), requestValidationInput); err != nil {
		var reqErr *openapi3filter.RequestError
		if errors.As(err, &reqErr) {
			return http.StatusBadRequest, err
		}
		var seqErr *openapi3filter.SecurityRequirementsError
		if errors.As(err, &seqErr) {
			return http.StatusUnauthorized, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// InvalidAPIEndpointHandler returns ErrInvalidAPIEndpoint, and is currently being used to ensure
// that routes under the pattern it is used with in chi.Router.Mount (i.e. /api/v1) are
// not accessible.
func InvalidAPIEndpointHandler(w http.ResponseWriter, r *http.Request) {
	writeError(w, r, http.StatusInternalServerError, ErrInvalidAPIEndpoint)
}
