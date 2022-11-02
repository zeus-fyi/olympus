package s3

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ConnectS3Session(ctx context.Context) (*s3.Client, error) {
	spacesKey := os.Getenv("SPACES_KEY")
	spacesSecret := os.Getenv("SPACES_SECRET")

	creds := credentials.NewStaticCredentialsProvider(spacesKey, spacesSecret, "")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: "https://sfo3.digitaloceanspaces.com",
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(creds),
		config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		return nil, err
	}
	// Create an Amazon S3 service client
	awsS3Client := s3.NewFromConfig(cfg)

	return awsS3Client, err
}
