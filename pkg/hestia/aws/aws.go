package hestia_eks_aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/rs/zerolog/log"
)

const (
	UsWest1           = "us-west-1"
	AwsUsWest1Context = "zeus-us-west-1"
	AwsEksRole        = "arn:aws:iam::480391564655:role/AWS-EKS-Role"
	AwsSubnetIDWest1A = "subnet-024b0ed90f92a7240"
	AwsSubnetIDWest1B = "subnet-086538f99779ddbc5"
)

var UsWestSubnetIDs = []string{AwsSubnetIDWest1A, AwsSubnetIDWest1B}

type AwsEKS struct {
	*eks.Client
}

func (a *AwsEKS) GetFullContextName(alias string) string {
	return fmt.Sprintf("arn:aws:eks:us-west-1:480391564655:cluster/%s", alias)
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

func (a *AwsEKS) AddNodeGroup(ctx context.Context, ngReq *eks.CreateNodegroupInput) (*eks.CreateNodegroupOutput, error) {
	ngReq.NodeRole = aws.String(AwsEksRole)
	ngReq.Subnets = []string{AwsSubnetIDWest1A, AwsSubnetIDWest1B}
	ngr, err := a.CreateNodegroup(ctx, ngReq)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return ngr, err
	}
	return ngr, err
}
