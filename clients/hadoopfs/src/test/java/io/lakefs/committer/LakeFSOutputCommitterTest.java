package io.lakefs;

import io.lakefs.clients.api.*;
import io.lakefs.clients.api.model.*;
import io.lakefs.clients.api.model.ObjectStats.PathTypeEnum;
import io.lakefs.utils.ObjectLocation;
import io.lakefs.LakeFSClient;

import com.amazonaws.ClientConfiguration;
import com.amazonaws.auth.AWSCredentials;
import com.amazonaws.auth.BasicAWSCredentials;
import com.amazonaws.services.s3.AmazonS3;
import com.amazonaws.services.s3.AmazonS3Client;
import com.amazonaws.services.s3.S3ClientOptions;
import com.amazonaws.services.s3.model.*;
import com.aventrix.jnanoid.jnanoid.NanoIdUtils;
import com.google.common.collect.ImmutableList;
import com.google.common.collect.Lists;
import org.apache.commons.io.IOUtils;
import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FileAlreadyExistsException;
import org.apache.hadoop.fs.FileStatus;
import org.apache.hadoop.fs.LocatedFileStatus;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.fs.FileSystem;
import org.apache.http.HttpStatus;

import org.apache.spark.SparkConf;
import org.apache.spark.sql.SparkSession;
import org.apache.spark.sql.Row;
import org.apache.spark.sql.RowFactory;
import org.apache.spark.sql.Dataset;
import org.apache.spark.sql.types.*;

import org.junit.Assert;
import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.After;
import org.junit.AfterClass;
import org.junit.Rule;
import org.junit.Test;
import org.junit.Ignore;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testcontainers.containers.*;
import org.testcontainers.containers.GenericContainer;
import org.testcontainers.containers.Network;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.lifecycle.Startables;
import org.testcontainers.containers.wait.strategy.Wait;
import org.testcontainers.utility.DockerImageName;

import java.io.*;
import java.net.URI;
import java.net.URISyntaxException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

public class LakeFSOutputCommitterTest {
    Logger logger = LoggerFactory.getLogger("LakeFSOutputCommitterTest");

    protected Configuration conf;

    protected LakeFSClient lfsClient;
    protected RepositoriesApi repositoriesApi;

    protected AmazonS3 s3Client;

    protected String s3Base;
    protected String s3Bucket;

    protected SparkSession spark;

    private static final DockerImageName LAKEFS_IMAGE = DockerImageName.parse("treeverse/lakefs:0.86.0");
    protected static final String LAKEFS_ACCESS_KEY = "AKIAIOSFODNN7EXAMPLE";
    protected static final String LAKEFS_SECRET_KEY = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY";

    private static final DockerImageName MINIO = DockerImageName.parse("minio/minio:RELEASE.2021-06-07T21-40-51Z");
    protected static final String S3_ACCESS_KEY_ID = "AKIArootkey";
    protected static final String S3_SECRET_ACCESS_KEY = "secret/minio/key=";

    Network network = Network.newNetwork();

    public final GenericContainer<?> s3 = new GenericContainer<>(MINIO.toString()).
            withCommand("minio", "server", "/data").
            withEnv("MINIO_ROOT_USER", S3_ACCESS_KEY_ID).
            withEnv("MINIO_ROOT_PASSWORD", S3_SECRET_ACCESS_KEY).
            withEnv("MINIO_DOMAIN", "s3.local.lakefs.io").
            withEnv("MINIO_UPDATE", "off").
            withNetwork(network).
            withNetworkAliases("minio").
            withExposedPorts(9000);

    public final PostgreSQLContainer<?> postgresDB = new PostgreSQLContainer<>("postgres:11")
            .withExposedPorts(5432)
            .withDatabaseName("postgres")
            .withUsername("lakefs")
            .withPassword("lakefs")
            .withNetwork(network)
            .withNetworkAliases("postgres");

    public final GenericContainer<?> lakeFS = new GenericContainer<>(LAKEFS_IMAGE.toString()).
            withEnv("LAKEFS_DATABASE_POSTGRES_CONNECTION_STRING", "postgres://lakefs:lakefs@postgres/postgres?sslmode=disable").
            withEnv("LAKEFS_DATABASE_TYPE", "postgres").
            withEnv("LAKEFS_BLOCKSTORE_TYPE", "s3").
            withEnv("LAKEFS_BLOCKSTORE_S3_FORCE_PATH_STYLE", "true").
//            withEnv("LAKEFS_BLOCKSTORE_S3_ENDPOINT", "http://s3.local.lakefs.io:9000").
            withEnv("LAKEFS_BLOCKSTORE_S3_ENDPOINT", "http://minio:9000").
            withEnv("LAKEFS_BLOCKSTORE_S3_CREDENTIALS_ACCESS_KEY_ID", S3_ACCESS_KEY_ID).
            withEnv("LAKEFS_BLOCKSTORE_S3_CREDENTIALS_SECRET_ACCESS_KEY", S3_SECRET_ACCESS_KEY).
            withEnv("LAKEFS_AUTH_ENCRYPT_SECRET_KEY", "some random secret string").
            withEnv("LAKEFS_STATS_ENABLED", "false").
            withEnv("LAKEFS_LOGGING_LEVEL", "DEBUG").
            waitingFor(Wait.forHttp("/")).
            withNetwork(network).
            withNetworkAliases("lakefs").
            withExposedPorts(8000).
            dependsOn(postgresDB, s3).
            withCommand("run");


    @After
    public void closeSparkSession(){
        spark.close();
    }


    @Before
    public void setUp() throws Exception {
        Startables.deepStart(s3, postgresDB, lakeFS).join();
//        Thread.sleep(20000);
        Container.ExecResult setupResult = lakeFS.execInContainer("/app/lakefs", "setup", "--user-name=tester", "--access-key-id", LAKEFS_ACCESS_KEY, "--secret-access-key", LAKEFS_SECRET_KEY);
        String setupstdout = setupResult.getStdout();
        int setupexitCode = setupResult.getExitCode();
        System.out.println("setupstdout is " + setupstdout);
        System.out.println("setupexitCode is " + setupexitCode);


        String lakeFSAddress = String.format("http://%s:%s/api/v1", lakeFS.getHost(), lakeFS.getMappedPort(8000));
        System.out.println("lakeFSAddress is " + lakeFSAddress);
        String s3Endpoint = String.format("http://%s:%s", s3.getHost(), s3.getMappedPort(9000));
        System.out.println("s3 is " + s3Endpoint);


        AWSCredentials creds = new BasicAWSCredentials(S3_ACCESS_KEY_ID, S3_SECRET_ACCESS_KEY);

        ClientConfiguration clientConfiguration = new ClientConfiguration()
                .withSignerOverride("AWSS3V4SignerType");
//        String s3Endpoint = String.format("http://s3.local.lakefs.io:%d", s3.getMappedPort(9000));

        s3Client = new AmazonS3Client(creds, clientConfiguration);

        S3ClientOptions s3ClientOptions = new S3ClientOptions()
            .withPathStyleAccess(true);
        s3Client.setS3ClientOptions(s3ClientOptions);
        s3Client.setEndpoint(s3Endpoint);

        s3Bucket = "example";
        s3Base = String.format("s3://%s", s3Bucket);
        CreateBucketRequest cbr = new CreateBucketRequest(s3Bucket);
        s3Client.createBucket(cbr);

        conf = new Configuration(false);
        conf.set("fs.lakefs.impl", "io.lakefs.LakeFSFileSystem");
        conf.set("fs.lakefs.access.key", LAKEFS_ACCESS_KEY);
        conf.set("fs.lakefs.secret.key", LAKEFS_SECRET_KEY);
        conf.set("fs.lakefs.endpoint", lakeFSAddress);
        conf.set("fs.s3a.impl", "org.apache.hadoop.fs.s3a.S3AFileSystem");
        conf.set(org.apache.hadoop.fs.s3a.Constants.ACCESS_KEY, S3_ACCESS_KEY_ID);
        conf.set(org.apache.hadoop.fs.s3a.Constants.SECRET_KEY, S3_SECRET_ACCESS_KEY);
        conf.set(org.apache.hadoop.fs.s3a.Constants.ENDPOINT, s3Endpoint);

        lfsClient = new LakeFSClient("lakefs://", conf);
        repositoriesApi = lfsClient.getRepositories();

        try {
            RepositoryCreation rc = new RepositoryCreation().name("example_repo").storageNamespace(s3Base);
            System.out.println("rc is " + rc);
            repositoriesApi.createRepository(rc, false);
        } catch (ApiException e) {
            System.out.println("error is " + e);
            System.out.println("lakeFS logs " + lakeFS.getLogs());
            System.out.println("Done logs ");
            System.out.println("s3 logs " + s3.getLogs());
            System.out.println("Done logs s3 ");
            System.out.println("going to sleep ");
            Thread.sleep(300000);
        }


        SparkConf sparkConf = new SparkConf().setMaster("local").setAppName("Spark test")
                .set("spark.sql.shuffle.partitions", "17")
                .set("spark.hadoop.mapred.output.committer.class","io.lakefs.committer.LakeFSOutputCommitter")
                .set("spark.hadoop.fs.lakefs.impl","io.lakefs.LakeFSFileSystem")
                .set("fs.lakefs.access.key", LAKEFS_ACCESS_KEY)
                .set("fs.lakefs.secret.key", LAKEFS_SECRET_KEY)
                .set("fs.lakefs.endpoint", lakeFSAddress)
                .set("fs.s3a.impl", "org.apache.hadoop.fs.s3a.S3AFileSystem")
                .set(org.apache.hadoop.fs.s3a.Constants.ACCESS_KEY, S3_ACCESS_KEY_ID)
                .set(org.apache.hadoop.fs.s3a.Constants.SECRET_KEY, S3_SECRET_ACCESS_KEY)
                .set(org.apache.hadoop.fs.s3a.Constants.ENDPOINT, s3Endpoint);
        this.spark = SparkSession.builder().config(sparkConf).getOrCreate();
    }

    @Ignore
    @Test
    public void testWriteTextfile() throws ApiException, IOException {
        String output = "lakefs://example_repo/main/some_df";

        StructType structType = new StructType().add("A", DataTypes.StringType, false);

        List<Row> nums = ImmutableList.of(
                RowFactory.create("value1"),
                RowFactory.create("value2"));
        Dataset<Row> df = spark.createDataFrame(nums, structType);

        ArrayList<String> failed = new ArrayList<String>();

        // TODO(Lynn): Add relevant file formats
        List<String> fileFormats = Arrays.asList("text");
        for (String fileFormat : fileFormats) {
            try {
                // save data in selected format
                String outputPath = String.format("%s.%s", output, fileFormat);
                logger.info(String.format("writing file in format %s to path %s", fileFormat, outputPath));
                df.write().format(fileFormat).save(outputPath);

                // read the data we just wrote
                logger.info(String.format("reading file in format %s", fileFormat));
                Dataset<Row> df_read = spark.read().format(fileFormat).load(outputPath);
                Dataset<Row> diff = df.exceptAll(df_read);
                Assert.assertTrue("all rows written were read", diff.isEmpty());
                diff = df_read.exceptAll(df);
                Assert.assertTrue("all rows read were written", diff.isEmpty());

            } catch (Exception e) {
                logger.error(String.format("Format %s unexpected exception: %s", fileFormat, e));
                failed.add(fileFormat);
            }
        }

        Assert.assertEquals("Failed list is not empty: %s" , 0, failed.size());
    }
}
