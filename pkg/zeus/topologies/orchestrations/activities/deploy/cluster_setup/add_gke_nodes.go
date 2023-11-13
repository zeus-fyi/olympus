package deploy_topology_activities_create_setup

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	hestia_gcp "github.com/zeus-fyi/olympus/pkg/hestia/gcp"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"google.golang.org/api/container/v1"
)

func (c *CreateSetupTopologyActivities) GkeAddNodePoolToOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest, npStatus do_types.DigitalOceanNodePoolRequestStatus) error {
	err := hestia_compute_resources.AddGkeNodePoolResourcesToOrg(ctx, params.Ou.OrgID, params.Nodes.ResourceID, params.NodesQuantity, npStatus.NodePoolID, npStatus.ClusterID, params.FreeTrial)
	if err != nil {
		log.Err(err).Interface("nodes", params.Nodes).Msg("GkeAddNodePoolToOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) SelectGkeNodeResources(ctx context.Context, request base_deploy_params.DestroyResourcesRequest) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	log.Info().Interface("request", request).Msg("SelectGkeNodeResources")
	nps, err := hestia_compute_resources.GkeSelectNodeResources(ctx, request.Ou.OrgID, request.OrgResourceIDs)
	if err != nil {
		log.Err(err).Interface("request", request).Msg("GkeSelectNodeResources: GkeSelectNodeResources error")
		return nps, err
	}
	return nps, err
}

func (c *CreateSetupTopologyActivities) GkeSelectFreeTrialNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	gkeNps, err := hestia_compute_resources.GkeSelectFreeTrialNodes(ctx, orgID)
	if err != nil {
		log.Err(err).Int("orgID", orgID).Msg("GkeSelectFreeTrialNodes: GkeSelectFreeTrialNodes error")
		return gkeNps, err
	}
	return gkeNps, err
}

func (c *CreateSetupTopologyActivities) GkeRemoveNodePoolRequest(ctx context.Context, nodePool do_types.DigitalOceanNodePoolRequestStatus) error {
	log.Info().Interface("nodePool", nodePool).Msg("GkeRemoveNodePoolRequest")
	ci := hestia_gcp.GcpClusterInfo{
		ClusterName: nodePool.ClusterID,
		ProjectID:   "zeusfyi",
		Zone:        "us-central1-a",
	}
	ni := hestia_gcp.GkeNodePoolInfo{
		Name: nodePool.NodePoolID,
	}
	_, err := api_auth_temporal.GCP.RemoveNodePool(ctx, ci, ni)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Info().Interface("nodePool", nodePool).Msg("GkeRemoveNodePoolRequest: node pool not found")
			return nil
		}
		log.Err(err).Interface("nodePool", nodePool).Msg("GkeRemoveNodePoolRequest error")
		return err
	}
	return err
}

func (c *CreateSetupTopologyActivities) GkeMakeNodePoolRequest(ctx context.Context, params base_deploy_params.ClusterSetupRequest) (do_types.DigitalOceanNodePoolRequestStatus, error) {
	log.Info().Interface("ou", params.Ou).Interface("nodes", params.Nodes).Msg("GkeMakeNodePoolRequest")
	labels := CreateBaseNodeLabels(params)

	tmp := strings.Split(params.CloudCtxNs.Namespace, "-")
	suffix := tmp[len(tmp)-1]
	tOrg := container.NodeTaint{
		Effect: "NO_SCHEDULE",
		Key:    fmt.Sprintf("org-%d", params.Ou.OrgID),
		Value:  fmt.Sprintf("org-%d", params.Ou.OrgID),
	}
	tApp := container.NodeTaint{
		Effect: "NO_SCHEDULE",
		Key:    "app",
		Value:  params.Cluster.ClusterName,
	}
	taints := []*container.NodeTaint{&tOrg}
	if params.AppTaint {
		taints = append(taints, &tApp)
	}
	// TODO remove hard code cluster info
	clusterID := "zeus-gcp-pilot-0"
	ci := hestia_gcp.GcpClusterInfo{
		ClusterName: clusterID,
		ProjectID:   "zeusfyi",
		Zone:        "us-central1-a",
	}
	name := strings.ToLower(fmt.Sprintf("gcp-%d-%s", params.Ou.OrgID, suffix))
	if len(name) > 39 {
		name = name[:39]
	}

	ni := hestia_gcp.GkeNodePoolInfo{
		Name:             name,
		MachineType:      params.Nodes.Slug,
		InitialNodeCount: int64(params.NodesQuantity),
	}

	node, err := api_auth_temporal.GCP.AddNodePool(ctx, ci, ni, taints, labels)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Info().Interface("nodeGroup", ni.Name).Msg("GkeMakeNodePoolRequest already exists")
			return do_types.DigitalOceanNodePoolRequestStatus{
				ClusterID:  clusterID,
				NodePoolID: ni.Name,
			}, nil
		}
		log.Err(err).Interface("node", node).Interface("nodes", params.Nodes).Msg("GkeMakeNodePoolRequest error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	//fmt.Println(node)
	return do_types.DigitalOceanNodePoolRequestStatus{
		ClusterID:  clusterID,
		NodePoolID: ni.Name,
	}, err
}
