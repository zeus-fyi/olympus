package artemis_orchestration_auth

import (
	"context"

	hestia_digitalocean "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean"
	hestia_gcp "github.com/zeus-fyi/olympus/pkg/hestia/gcp"
)

var (
	Bearer       string
	DigitalOcean hestia_digitalocean.DigitalOcean
	GCP          hestia_gcp.GcpClient
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
