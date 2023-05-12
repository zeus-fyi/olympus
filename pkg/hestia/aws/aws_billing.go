package hestia_eks_aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/rs/zerolog/log"
)

type AwsPricing struct {
	*pricing.Client
}

func InitPricingClient(ctx context.Context, accessCred EksCredentials) (AwsPricing, error) {
	creds := credentials.NewStaticCredentialsProvider(accessCred.AccessKey, accessCred.AccessSecret, "")
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(creds),
		config.WithRegion(accessCred.Region),
	)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return AwsPricing{}, err
	}
	return AwsPricing{pricing.NewFromConfig(cfg)}, nil
}
