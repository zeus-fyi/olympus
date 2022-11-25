package s3base

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	AwsS3Client    *s3.Client
	SpacesEndpoint string // TODO should set the aws url endpoint
	spacesKey      string
	spacesSecret   string
}

func NewS3ClientBase() S3Client {
	return S3Client{}
}

func NewS3ClientWithEndpoint(endpoint string) S3Client {
	return S3Client{SpacesEndpoint: endpoint}
}

func NewConnS3ClientWithStaticCreds(ctx context.Context, key, secret string) (S3Client, error) {
	s3client := S3Client{spacesKey: key, spacesSecret: secret}
	err := s3client.ConnectS3SpacesDO(ctx)
	return s3client, err
}

func (s *S3Client) ConnectS3SpacesDO(ctx context.Context) error {
	if len(s.spacesKey) <= 0 || len(s.spacesSecret) <= 0 {
		s.spacesKey = os.Getenv("SPACES_KEY")
		s.spacesSecret = os.Getenv("SPACES_SECRET")
	}
	creds := credentials.NewStaticCredentialsProvider(s.spacesKey, s.spacesSecret, "")
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: "https://sfo3.digitaloceanspaces.com",
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(creds),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRetryMaxAttempts(100))
	if err != nil {
		return err
	}
	// Create an Amazon S3 service client
	s.AwsS3Client = s3.NewFromConfig(cfg)

	return err
}
