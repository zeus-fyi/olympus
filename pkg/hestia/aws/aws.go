package hestia_eks_aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/rs/zerolog/log"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

const (
	UsWest1                  = "us-west-1"
	AwsUsWest1Context        = "zeus-us-west-1"
	AwsEksRole               = "arn:aws:iam::480391564655:role/AWS-EKS-Role"
	AwsUsWestSecurityGroupID = "sg-0f62afc9340e7df70"
	AwsSubnetIDWest1A        = "subnet-024b0ed90f92a7240"
	AwsSubnetIDWest1B        = "subnet-086538f99779ddbc5"
)

var UsWestSubnetIDs = []string{AwsSubnetIDWest1A, AwsSubnetIDWest1B}

type AwsEKS struct {
	*eks.Client
}

func (a *AwsEKS) GetFullContextName(alias string) string {
	return fmt.Sprintf("arn:aws:eks:us-west-1:480391564655:cluster/%s", alias)
}

func InitAwsEKS(ctx context.Context, accessCred aegis_aws_auth.AuthAWS) (AwsEKS, error) {
	creds := credentials.NewStaticCredentialsProvider(accessCred.AccessKey, accessCred.SecretKey, "")
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

func (a *AwsEKS) RemoveNodeGroup(ctx context.Context, ngReq *eks.DeleteNodegroupInput) (*eks.DeleteNodegroupOutput, error) {
	ngr, err := a.DeleteNodegroup(ctx, ngReq)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return ngr, err
	}
	return ngr, err
}
