package hestia_eks_aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"
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
	*eks.Client         // Embed the client
	Arn         *string `json:"arn"`
	Account     *string `json:"account"`
	Username    string  `json:"username"`
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
		log.Err(err).Msg("InitAwsEKS")
		return AwsEKS{}, err
	}
	// Create an Amazon STS client from just a session.
	stsClient := sts.NewFromConfig(cfg)

	// Call to get the caller identity
	result, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return AwsEKS{}, err
	}
	var username string
	if result.Arn != nil {
		username = strings.Split(*result.Arn, "/")[len(strings.Split(*result.Arn, "/"))-1]
	}
	return AwsEKS{Arn: result.Arn, Account: result.Account, Client: eks.NewFromConfig(cfg), Username: username}, nil
}

func (a *AwsEKS) AddNodeGroup(ctx context.Context, ngReq *eks.CreateNodegroupInput) (*eks.CreateNodegroupOutput, error) {
	ngr, err := a.CreateNodegroup(ctx, ngReq)
	if err != nil {
		log.Err(err).Msg("AddNodeGroup")
		return ngr, err
	}
	return ngr, err
}

func (a *AwsEKS) RemoveNodeGroup(ctx context.Context, ngReq *eks.DeleteNodegroupInput) (*eks.DeleteNodegroupOutput, error) {
	ngr, err := a.DeleteNodegroup(ctx, ngReq)
	if err != nil {
		log.Err(err).Msg("RemoveNodeGroup")
		return ngr, err
	}
	return ngr, err
}
