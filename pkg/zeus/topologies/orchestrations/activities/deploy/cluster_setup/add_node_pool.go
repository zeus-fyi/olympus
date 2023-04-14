package deploy_topology_activities_create_setup

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
)

func (c *CreateSetupTopologyActivities) AddNodePoolToOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest, npStatus do_types.DigitalOceanNodePoolRequestStatus) error {
	err := hestia_compute_resources.AddDigitalOceanNodePoolResourcesToOrg(ctx, params.Ou.OrgID, params.Nodes.ResourceID, params.NodesQuantity, npStatus.NodePoolID, npStatus.ClusterID, params.FreeTrial)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("nodes", params.Nodes).Msg("AddNodePoolToOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) MakeNodePoolRequest(ctx context.Context, params base_deploy_params.ClusterSetupRequest) (do_types.DigitalOceanNodePoolRequestStatus, error) {
	taint := godo.Taint{
		Key:    fmt.Sprintf("org-%d", params.Ou.OrgID),
		Value:  fmt.Sprintf("org-%d", params.Ou.OrgID),
		Effect: "NoSchedule",
	}
	label := make(map[string]string)
	label["org"] = fmt.Sprintf("%d", params.Ou.OrgID)
	label["app"] = params.Cluster.ClusterName
	suffix := strings.Split(params.ClusterID.String(), "-")[0]
	nodesReq := &godo.KubernetesNodePoolCreateRequest{
		Name:   fmt.Sprintf("nodepool-%d-%s", params.Ou.OrgID, suffix),
		Size:   params.Nodes.Slug,
		Count:  int(params.NodesQuantity),
		Tags:   nil,
		Labels: label,
		Taints: []godo.Taint{taint},
	}
	// TODO remove hard code cluster id
	clusterID := "0de1ee8e-7b90-45ea-b966-e2d2b7976cf9"
	node, err := api_auth_temporal.DigitalOcean.CreateNodePool(ctx, clusterID, nodesReq)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("nodes", params.Nodes).Msg("CreateNodePool error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}

	return do_types.DigitalOceanNodePoolRequestStatus{
		ClusterID:  clusterID,
		NodePoolID: node.ID,
	}, nil
}

// clusterID := "0de1ee8e-7b90-45ea-b966-e2d2b7976cf9"
func (c *CreateSetupTopologyActivities) RemoveNodePoolRequest(ctx context.Context, nodePool do_types.DigitalOceanNodePoolRequestStatus) error {
	log.Ctx(ctx).Info().Interface("nodePool", nodePool).Msg("RemoveNodePoolRequest")
	err := api_auth_temporal.DigitalOcean.RemoveNodePool(ctx, nodePool.ClusterID, nodePool.NodePoolID)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("nodePool", nodePool).Msg("RemoveNodePool error")
		return err
	}
	return err
}

func (c *CreateSetupTopologyActivities) RemoveFreeTrialOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest) error {
	err := hestia_compute_resources.RemoveFreeTrialOrgResources(ctx, params.Ou.OrgID)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", params.Ou).Msg("RemoveFreeTrialOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) UpdateFreeTrialOrgResourcesToPaid(ctx context.Context, params base_deploy_params.ClusterSetupRequest) error {
	err := hestia_compute_resources.UpdateFreeTrialOrgResourcesToPaid(ctx, params.Ou.OrgID)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", params.Ou).Msg("RemoveFreeTrialOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) SelectFreeTrialNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	nps, err := hestia_compute_resources.SelectFreeTrialDigitalOceanNodes(ctx, orgID)
	if err != nil {
		log.Ctx(ctx).Err(err).Int("orgID", orgID).Msg("SelectFreeTrialNodes: SelectFreeTrialDigitalOceanNodes error")
		return nps, err
	}
	return nps, err
}
