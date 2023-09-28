package hestia_digitalocean

import (
	"context"
	"net/http"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
)

func AddDoNvmeLabels(labels map[string]string) map[string]string {
	labels["fast-disk-node"] = "pv-raid"
	return labels
}

func (d *DigitalOcean) CreateNodePool(ctx context.Context, context string, nodesReq *godo.KubernetesNodePoolCreateRequest) (*godo.KubernetesNodePool, error) {
	nodePool, _, err := d.Kubernetes.CreateNodePool(ctx, context, nodesReq)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create node pool")
		return nil, err
	}
	return nodePool, err
}

func (d *DigitalOcean) RemoveNodePool(ctx context.Context, context, poolID string) error {
	r, err := d.Kubernetes.DeleteNodePool(ctx, context, poolID)
	if r.StatusCode == http.StatusNotFound {
		log.Ctx(ctx).Error().Err(err).Msg("node pool doesn't exist, it may have been removed already")
		return nil
	}
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create node pool")
		return err
	}
	return err
}
