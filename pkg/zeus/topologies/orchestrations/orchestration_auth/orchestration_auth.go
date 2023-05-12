package artemis_orchestration_auth

import (
	"context"

	hestia_eks_aws "github.com/zeus-fyi/olympus/pkg/hestia/aws"
	hestia_digitalocean "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean"
	hestia_gcp "github.com/zeus-fyi/olympus/pkg/hestia/gcp"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

var (
	Bearer       string
	DigitalOcean hestia_digitalocean.DigitalOcean
	GCP          hestia_gcp.GcpClient
	Eks          hestia_eks_aws.AwsEKS
)

func InitOrchestrationDigitalOceanClient(ctx context.Context, bearer string) {
	Bearer = bearer
	DigitalOcean = hestia_digitalocean.InitDoClient(ctx, Bearer)
}

func InitOrchestrationGcpClient(ctx context.Context, authJsonBytes []byte) {
	g, err := hestia_gcp.InitGcpClient(ctx, authJsonBytes)
	if err != nil {
		panic(err)
	}
	GCP = g
}

func InitOrchestrationEksClient(ctx context.Context, accessCred aegis_aws_auth.AuthAWS) {
	eks, err := hestia_eks_aws.InitAwsEKS(ctx, accessCred)
	if err != nil {
		panic(err)
	}
	Eks = eks
}
