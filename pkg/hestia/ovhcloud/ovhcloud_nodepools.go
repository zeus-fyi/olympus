package hestia_ovhcloud

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (o *OvhCloud) CreateNodePool(ctx context.Context, nodesReq OvhNodePoolCreationRequest) (*OvhNodePoolCreationResponse, error) {
	resp := &OvhNodePoolCreationResponse{}
	err := o.CallAPIWithContext(ctx, "POST", nodesReq.GetEndpoint(), nodesReq, resp, true)
	if err != nil {
		log.Err(err).Msg("CreateNodePool: Ovh")
	}
	return resp, nil
}

// /cloud/project/{serviceName}/kube/{kubeId}/nodepool/{nodePoolId}

func (o *OvhCloud) RemoveNodePool(ctx context.Context, context, poolID string) error {
	return nil
}
