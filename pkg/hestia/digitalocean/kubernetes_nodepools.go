package hestia_digitalocean

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
)

func (d *DigitalOcean) CreateNodePool(ctx context.Context, context string, nodesReq *godo.KubernetesNodePoolCreateRequest) (*godo.KubernetesNodePool, error) {
	nodePool, _, err := d.Kubernetes.CreateNodePool(ctx, context, nodesReq)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create node pool")
		return nil, err
	}
	return nodePool, err
}

func (d *DigitalOcean) RemoveNodePool(ctx context.Context, context, poolID string) error {
	_, err := d.Kubernetes.DeleteNodePool(ctx, context, poolID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create node pool")
		return err
	}
	return err
}
