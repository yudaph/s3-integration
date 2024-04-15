package s3client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"sync"
)

var (
	uploader *manager.Uploader
	once     sync.Once
)

func setup() {
	endpoint := "http://s3-test-endpoint.com"
	accessKey := "XVXQ8Q3QGQ9QXQ8QXQ8Q"
	secretKey := "/XQ8Q3QGQ9QXQ8QXQ8QXQ/Q3QGQ9QXQ8QXQ8QXQ8"

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               endpoint,
			HostnameImmutable: true,
			Source:            aws.EndpointSourceCustom,
		}, nil
	})
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(resolver),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatalf("failed to load config, %v", err)
	}

	uploader = manager.NewUploader(s3.NewFromConfig(awsConfig)) // Initialize S3 client
}

func UploadFile(ctx context.Context, bucket, filename string, file io.Reader) (string, error) {
	once.Do(setup)
	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file,
		ACL:    "public-read",
	})
	if err != nil {
		log.Fatalf("failed to upload file, %v", err)
	}

	return result.Location, nil
}
