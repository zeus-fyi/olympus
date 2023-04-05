package hestia_digitalocean

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
)

func (d *DigitalOcean) AddToNodePool(ctx context.Context) (*godo.KubernetesNodePool, error) {
	createRequest := &godo.KubernetesNodePoolCreateRequest{
		Name:   "pool-02",
		Size:   "s-2vcpu-4gb",
		Count:  1,
		Tags:   []string{"web"},
		Labels: map[string]string{"service": "web", "priority": "high"},
	}
	nodePool, _, err := d.Kubernetes.CreateNodePool(ctx, "cluster-uuid", createRequest)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create node pool")
	}
	return nodePool, err
}
