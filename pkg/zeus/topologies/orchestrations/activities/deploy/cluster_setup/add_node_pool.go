package deploy_topology_activities_create_setup

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
)

func (c *CreateSetupTopologyActivities) AddNodePoolToOrgResources(ctx context.Context, orgID int, nodes hestia_autogen_bases.Nodes, quantity float64) error {
	err := hestia_compute_resources.AddResourcesToOrg(ctx, orgID, nodes.ResourceID, quantity)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("nodes", nodes).Msg("AddNodePoolToOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) MakeNodePoolRequest(ctx context.Context, clusterID string, nodes hestia_autogen_bases.Nodes, quantity float64) error {
	nodesReq := &godo.KubernetesNodePoolCreateRequest{
		Name:   clusterID + "-node-pool",
		Size:   "nodes.Slug",
		Count:  int(quantity),
		Tags:   nil,
		Labels: nil,
		Taints: nil, // TODO: add taints
	}
	_, err := api_auth_temporal.DigitalOcean.AddToNodePool(ctx, clusterID, nodesReq)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("nodes", nodes).Msg("MakeNodePoolRequest error")
		return err
	}
	return nil
}
