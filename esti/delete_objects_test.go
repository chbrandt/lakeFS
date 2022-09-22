package esti

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/treeverse/lakefs/pkg/api"
	"github.com/treeverse/lakefs/pkg/testutil"
)

func testDeleteObjects(t *testing.T, svc *s3.S3) {
	ctx, _, repo := setupTest(t)
	defer tearDownTest(repo)
	const numOfObjects = 10

	identifiers := make([]*s3.ObjectIdentifier, 0, numOfObjects)

	for i := 1; i <= numOfObjects; i++ {
		file := strconv.Itoa(i) + ".txt"
		identifiers = append(identifiers, &s3.ObjectIdentifier{
			Key: aws.String(mainBranch + "/" + file),
		})
		_, _ = uploadFileRandomData(ctx, t, repo, mainBranch, file, false)
	}

	listOut, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(repo),
		Prefix: aws.String(mainBranch + "/"),
	})

	assert.NoError(t, err)
	assert.Len(t, listOut.Contents, numOfObjects)

	deleteOut, err := svc.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(repo),
		Delete: &s3.Delete{
			Objects: identifiers,
		},
	})

	assert.NoError(t, err)
	assert.Len(t, deleteOut.Deleted, numOfObjects)

	listOut, err = svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(repo),
		Prefix: aws.String(mainBranch + "/"),
	})

	assert.NoError(t, err)
	assert.Len(t, listOut.Contents, 0)
}

func TestDeleteObjects(t *testing.T) {
	testDeleteObjects(t, pathStyleSvc)
	if !skipS3HostStyleTests {
		testDeleteObjects(t, hostStyleSvc)
	}
}

func testDeleteObjectsViewer(t *testing.T, useHostBaseClient bool, svc *s3.S3) {
	ctx, _, repo := setupTest(t)
	defer tearDownTest(repo)

	// setup data
	const filename = "delete-me"
	_, _ = uploadFileRandomData(ctx, t, repo, mainBranch, filename, false)

	// setup user with only view rights - create user, add to group, generate credentials
	uid := "del-viewer"
	resCreateUser, err := client.CreateUserWithResponse(ctx, api.CreateUserJSONRequestBody{
		Id: uid,
	})
	require.NoError(t, err, "Admin failed while creating user")
	require.Equal(t, http.StatusCreated, resCreateUser.StatusCode(), "Admin unexpectedly failed to create user")

	resAssociateUser, err := client.AddGroupMembershipWithResponse(ctx, "Viewers", "del-viewer")
	require.NoError(t, err, "Failed to add user to Viewers group")
	require.Equal(t, http.StatusCreated, resAssociateUser.StatusCode(), "AddGroupMembershipWithResponse unexpectedly status code")

	resCreateCreds, err := client.CreateCredentialsWithResponse(ctx, "del-viewer")
	require.NoError(t, err, "Failed to create credentials")
	require.Equal(t, http.StatusCreated, resCreateCreds.StatusCode(), "CreateCredentials unexpectedly status code")

	// client with viewer user credentials
	creds := resCreateCreds.JSON201
	svcViewer := testutil.SetupTestS3Client(creds.AccessKeyId, creds.SecretAccessKey, s3Endpoint, useHostBaseClient)

	// delete objects using viewer
	deleteOut, err := svcViewer.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(repo),
		Delete: &s3.Delete{
			Objects: []*s3.ObjectIdentifier{{Key: api.StringPtr(mainBranch + "/" + filename)}},
		},
	})
	// make sure we got an error we fail to delete the file
	assert.NoError(t, err)
	assert.Len(t, deleteOut.Errors, 1, "error we fail to delete")
	assert.Len(t, deleteOut.Deleted, 0, "no file should be deleted")

	// verify that viewer can't delete the file
	listOut, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(repo),
		Prefix: aws.String(mainBranch + "/"),
	})
	assert.NoError(t, err)
	assert.Len(t, listOut.Contents, 1, "list should find 'delete-me' file")
	assert.Equal(t, aws.StringValue(listOut.Contents[0].Key), mainBranch+"/"+filename)
}

func TestDeleteObjects_Viewer(t *testing.T) {
	t.SkipNow()
	testDeleteObjectsViewer(t, false, pathStyleSvc)
	if !skipS3HostStyleTests {
		testDeleteObjectsViewer(t, true, hostStyleSvc)
	}
}
