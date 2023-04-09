package deploy_topology_activities_create_setup

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	deploy_workflow_cluster_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create_setup"
)

func (c *CreateSetupTopologyActivities) AddNodePoolToOrgResources(ctx context.Context, params deploy_workflow_cluster_setup.ClusterSetupRequest) error {
	err := hestia_compute_resources.AddResourcesToOrg(ctx, params.Ou.OrgID, params.Nodes.ResourceID, params.NodesQuantity)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("nodes", params.Nodes).Msg("AddNodePoolToOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) MakeNodePoolRequest(ctx context.Context, params deploy_workflow_cluster_setup.ClusterSetupRequest) error {
	taint := godo.Taint{
		Key:   fmt.Sprintf("%d", params.Ou.OrgID),
		Value: params.ClusterID.String(),
	}
	nodesReq := &godo.KubernetesNodePoolCreateRequest{
		Name:   params.ClusterID.String() + "-node-pool",
		Size:   params.Nodes.Slug,
		Count:  int(params.NodesQuantity),
		Tags:   nil,
		Labels: nil,
		Taints: []godo.Taint{taint},
	}
	_, err := api_auth_temporal.DigitalOcean.AddToNodePool(ctx, params.ClusterID.String(), nodesReq)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("nodes", params.Nodes).Msg("MakeNodePoolRequest error")
		return err
	}
	return nil
}
