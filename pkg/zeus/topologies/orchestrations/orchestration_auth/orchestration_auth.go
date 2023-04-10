package artemis_orchestration_auth

import (
	"context"

	hestia_digitalocean "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean"
)

var (
	Bearer       string
	DigitalOcean hestia_digitalocean.DigitalOcean
)

func InitOrchestrationDigitalOceanClient(ctx context.Context, bearer string) {
	Bearer = bearer
	DigitalOcean = hestia_digitalocean.InitDoClient(ctx, Bearer)
}
