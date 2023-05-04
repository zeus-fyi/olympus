package deploy_topology_activities_create_setup

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	hestia_gcp "github.com/zeus-fyi/olympus/pkg/hestia/gcp"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"google.golang.org/api/container/v1"
)

func (c *CreateSetupTopologyActivities) AddNodePoolToOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest, npStatus do_types.DigitalOceanNodePoolRequestStatus) error {
	err := hestia_compute_resources.AddDigitalOceanNodePoolResourcesToOrg(ctx, params.Ou.OrgID, params.Nodes.ResourceID, params.NodesQuantity, npStatus.NodePoolID, npStatus.ClusterID, params.FreeTrial)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("nodes", params.Nodes).Msg("AddNodePoolToOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) GkeAddNodePoolToOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest, npStatus do_types.DigitalOceanNodePoolRequestStatus) error {
	err := hestia_compute_resources.AddGkeNodePoolResourcesToOrg(ctx, params.Ou.OrgID, params.Nodes.ResourceID, params.NodesQuantity, npStatus.NodePoolID, npStatus.ClusterID, params.FreeTrial)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("nodes", params.Nodes).Msg("GkeAddNodePoolToOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) GkeMakeNodePoolRequest(ctx context.Context, params base_deploy_params.ClusterSetupRequest) (do_types.DigitalOceanNodePoolRequestStatus, error) {
	label := make(map[string]string)
	label["org"] = fmt.Sprintf("%d", params.Ou.OrgID)
	label["app"] = params.Cluster.ClusterName
	suffix := strings.Split(params.Namespace, "-")[0]
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
	taints := []*container.NodeTaint{&tOrg, &tApp}
	// TODO remove hard code cluster info
	clusterID := "zeus-gcp-pilot-0"
	ci := hestia_gcp.GcpClusterInfo{
		ClusterName: clusterID,
		ProjectID:   "zeusfyi",
		Zone:        "us-central1-a",
	}
	ni := hestia_gcp.GkeNodePoolInfo{
		Name:             fmt.Sprintf("nodepool-%d-%s", params.Ou.OrgID, suffix),
		MachineType:      params.Nodes.Slug,
		InitialNodeCount: int64(params.NodesQuantity),
	}
	node, err := api_auth_temporal.GCP.AddNodePool(ctx, ci, ni, taints)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("nodes", params.Nodes).Msg("CreateNodePool error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	fmt.Println(node)
	return do_types.DigitalOceanNodePoolRequestStatus{
		ClusterID:  clusterID,
		NodePoolID: ni.Name,
	}, err
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
	suffix := strings.Split(params.Namespace, "-")[0]
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

func (c *CreateSetupTopologyActivities) SelectGkeNodeResources(ctx context.Context, request base_deploy_params.DestroyResourcesRequest) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	log.Ctx(ctx).Info().Interface("request", request).Msg("SelectNodeResources")
	nps, err := hestia_compute_resources.GkeSelectNodeResources(ctx, request.Ou.OrgID, request.OrgResourceIDs)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("request", request).Msg("GkeSelectNodeResources: GkeSelectNodeResources error")
		return nps, err
	}
	return nps, err
}

func (c *CreateSetupTopologyActivities) SelectNodeResources(ctx context.Context, request base_deploy_params.DestroyResourcesRequest) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	log.Ctx(ctx).Info().Interface("request", request).Msg("SelectNodeResources")
	nps, err := hestia_compute_resources.SelectNodeResources(ctx, request.Ou.OrgID, request.OrgResourceIDs)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("request", request).Msg("SelectNodeResources: SelectNodeResources error")
		return nps, err
	}
	return nps, err
}

func (c *CreateSetupTopologyActivities) EndResourceService(ctx context.Context, request base_deploy_params.DestroyResourcesRequest) error {
	log.Ctx(ctx).Info().Interface("request", request).Msg("EndResourceService")
	err := hestia_compute_resources.UpdateEndServiceOrgResources(ctx, request.Ou.OrgID, request.OrgResourceIDs)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("request", request).Msg("EndResourceService: UpdateEndServiceOrgResources error")
		return err
	}
	return err
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

func (c *CreateSetupTopologyActivities) GkeRemoveNodePoolRequest(ctx context.Context, nodePool do_types.DigitalOceanNodePoolRequestStatus) error {
	log.Ctx(ctx).Info().Interface("nodePool", nodePool).Msg("RemoveNodePoolRequest")
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
		log.Ctx(ctx).Err(err).Interface("nodePool", nodePool).Msg("GkeRemoveNodePoolRequest error")
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

func (c *CreateSetupTopologyActivities) GkeSelectFreeTrialNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	gkeNps, err := hestia_compute_resources.GkeSelectFreeTrialNodes(ctx, orgID)
	if err != nil {
		log.Ctx(ctx).Err(err).Int("orgID", orgID).Msg("GkeSelectFreeTrialNodes: GkeSelectFreeTrialNodes error")
		return gkeNps, err
	}
	return gkeNps, err
}
