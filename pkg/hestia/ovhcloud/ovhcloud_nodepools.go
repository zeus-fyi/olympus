package hestia_ovhcloud

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (o *OvhCloud) GetNodePool(ctx context.Context, nodesReq OvhNodePoolCreationRequest) (any, error) {
	var resp any
	err := o.GetWithContext(ctx, nodesReq.GetEndpoint(), resp)
	if err != nil {
		log.Err(err).Msg("CreateNodePool: Ovh")
		return nil, err
	}
	return resp, nil
}

func (o *OvhCloud) CreateNodePool(ctx context.Context, nodesReq OvhNodePoolCreationRequest) (*OvhNodePoolCreationResponse, error) {
	resp := &OvhNodePoolCreationResponse{}
	endpoint := nodesReq.GetEndpoint()
	err := o.PostWithContext(ctx, endpoint, nodesReq.ProjectKubeNodePoolCreation, resp)
	if err != nil {
		log.Err(err).Msg("CreateNodePool: Ovh")
		return nil, err
	}
	return resp, nil
}

func (o *OvhCloud) RemoveNodePool(ctx context.Context, context, poolID string) error {
	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", OvhServiceName, context, poolID)
	err := o.CallAPIWithContext(ctx, "DELETE", endpoint, nil, nil, true)
	if err != nil {
		log.Err(err).Msg("RemoveNodePool: Ovh")
		return err
	}
	return nil
}
