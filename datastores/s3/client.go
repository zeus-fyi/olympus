package s3base

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
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
	if err != nil {
		log.Err(err).Msg("NewConnS3ClientWithStaticCreds: ConnectS3SpacesDO failed")
		panic(err)
	}
	return s3client, err
}

func (s *S3Client) ConnectS3SpacesDO(ctx context.Context) error {
	if len(s.spacesKey) <= 0 || len(s.spacesSecret) <= 0 {
		log.Info().Msg("S3Client: ConnectS3SpacesDO had no provided param credentials, checking env vars")
		s.spacesKey = os.Getenv("SPACES_KEY")
		s.spacesSecret = os.Getenv("SPACES_SECRET")
	}
	if len(s.spacesKey) <= 0 || len(s.spacesSecret) <= 0 {
		panic("S3Client: ConnectS3SpacesDO had no provided param credentials, and no env vars")
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
		log.Err(err).Msg("ConnectS3SpacesDO: config.LoadDefaultConfig failed")
		panic(err)
	}
	// Create an Amazon S3 service client
	s.AwsS3Client = s3.NewFromConfig(cfg)
	return err
}
