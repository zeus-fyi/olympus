package hestia_eks_aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/rs/zerolog/log"
)

const UsWest1 = "us-west-1"

type AwsEKS struct {
	*eks.Client
}

type EksCredentials struct {
	Region       string
	AccessKey    string
	AccessSecret string
}

func InitAwsEKS(ctx context.Context, accessCred EksCredentials) (AwsEKS, error) {
	creds := credentials.NewStaticCredentialsProvider(accessCred.AccessKey, accessCred.AccessSecret, "")
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(creds),
		config.WithRegion(accessCred.Region),
	)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return AwsEKS{}, err
	}
	return AwsEKS{eks.NewFromConfig(cfg)}, nil
}
