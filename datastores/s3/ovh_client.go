package s3base

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

func NewOvhConnS3ClientWithStaticCreds(ctx context.Context, key, secret string) (S3Client, error) {
	s3client := S3Client{spacesKey: key, spacesSecret: secret}
	err := s3client.ConnectS3SpacesOvh(ctx)
	if err != nil {
		log.Err(err).Msg("NewConnS3ClientWithStaticCreds: ConnectS3SpacesDO failed")
		misc.DelayedPanic(err)
	}
	return s3client, err
}
func (s *S3Client) ConnectS3SpacesOvh(ctx context.Context) error {
	if len(s.spacesKey) <= 0 && len(s.spacesSecret) <= 0 {
		panic("S3Client: } had no provided param credentials, and no env vars")
	}
	creds := credentials.NewStaticCredentialsProvider(s.spacesKey, s.spacesSecret, "")
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: "https://s3.us-west-or.perf.cloud.ovh.us",
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-or"),
		config.WithCredentialsProvider(creds),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRetryMaxAttempts(100))
	if err != nil {
		log.Err(err).Msg("ConnectS3SpacesDO: config.LoadDefaultConfig failed")
		misc.DelayedPanic(err)
	}
	// Create an Amazon S3 service client
	s.AwsS3Client = s3.NewFromConfig(cfg)
	return err
}
