package hestia_ovhcloud

import (
	"context"
	"fmt"

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
	endpoint := fmt.Sprintf("/v1/cloud/project/%s/kube/nodepool/%s", context, poolID)
	err := o.CallAPIWithContext(ctx, "DELETE", endpoint, "", "", true)
	if err != nil {
		log.Err(err).Msg("RemoveNodePool: Ovh")
	}
	return nil
}
