package deploy_topology_activities_create_setup

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	hestia_digitalocean "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean"
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	hestia_ovhcloud "github.com/zeus-fyi/olympus/pkg/hestia/ovhcloud"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
)

func (c *CreateSetupTopologyActivities) AddNodePoolToOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest, npStatus do_types.DigitalOceanNodePoolRequestStatus) error {
	err := hestia_compute_resources.AddDigitalOceanNodePoolResourcesToOrg(ctx, params.Ou.OrgID, params.Nodes.ResourceID, params.NodesQuantity, npStatus.NodePoolID, npStatus.ClusterID, params.FreeTrial)
	if err != nil {
		log.Err(err).Interface("nodes", params.Nodes).Msg("AddNodePoolToOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) OvhAddNodePoolToOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest, npStatus do_types.DigitalOceanNodePoolRequestStatus) error {
	err := hestia_compute_resources.AddOvhNodePoolResourcesToOrg(ctx, params.Ou.OrgID, params.Nodes.ResourceID, params.NodesQuantity, npStatus.NodePoolID, npStatus.ClusterID, params.FreeTrial)
	if err != nil {
		log.Err(err).Interface("nodes", params.Nodes).Msg("OvhAddNodePoolToOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) OvhMakeNodePoolRequest(ctx context.Context, params base_deploy_params.ClusterSetupRequest) (do_types.DigitalOceanNodePoolRequestStatus, error) {
	kubeId := hestia_ovhcloud.OvhSharedKubeID
	switch params.Ou.UserID {
	case 7138958574876245565:
		if params.Ou.OrgID == 7138983863666903883 {
			kubeId = hestia_ovhcloud.OvhInternalKubeID
		}
	}
	autoscaleEnabled := false
	tmp := strings.Split(params.CloudCtxNs.Namespace, "-")
	suffix := tmp[len(tmp)-1]

	nodeGroupName := strings.ToLower(fmt.Sprintf("ovh-%d-%s", params.Ou.OrgID, suffix))
	if len(nodeGroupName) > 39 {
		nodeGroupName = nodeGroupName[:39]
	}
	labels := CreateBaseNodeLabels(params)

	taints := []hestia_ovhcloud.KubernetesTaint{
		{
			Key:    fmt.Sprintf("org-%d", params.Ou.OrgID),
			Value:  fmt.Sprintf("org-%d", params.Ou.OrgID),
			Effect: "NoSchedule",
		},
	}
	if params.AppTaint && params.Cluster.ClusterName != "" {
		taints = append(taints, hestia_ovhcloud.KubernetesTaint{
			Key:    "app",
			Value:  params.Cluster.ClusterName,
			Effect: "NoSchedule",
		})
	}
	npr := hestia_ovhcloud.OvhNodePoolCreationRequest{
		ServiceName: hestia_ovhcloud.OvhServiceName,
		KubeId:      kubeId,
		ProjectKubeNodePoolCreation: hestia_ovhcloud.ProjectKubeNodePoolCreation{
			AntiAffinity:  nil,
			Autoscale:     &autoscaleEnabled,
			Autoscaling:   nil,
			DesiredNodes:  int(params.NodesQuantity),
			FlavorName:    params.Nodes.Slug,
			MaxNodes:      int(params.NodesQuantity),
			MinNodes:      int(params.NodesQuantity),
			MonthlyBilled: nil,
			Name:          nodeGroupName,
			Template: &hestia_ovhcloud.NodeTemplate{
				Metadata: &hestia_ovhcloud.Metadata{
					Annotations: make(map[string]string),
					Finalizers:  []string{},
					Labels:      labels,
				},
				Spec: &hestia_ovhcloud.Spec{
					Taints:        taints,
					Unschedulable: false,
				},
			},
		},
	}
	resp, err := api_auth_temporal.OvhCloud.CreateNodePool(ctx, npr)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Info().Msg("Node pool already exists")
			return do_types.DigitalOceanNodePoolRequestStatus{
				ClusterID:  kubeId,
				NodePoolID: resp.Id,
			}, nil
		}
		log.Err(err).Interface("nodes", params.Nodes).Msg("OvhMakeNodePoolRequest error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	return do_types.DigitalOceanNodePoolRequestStatus{
		ClusterID:  kubeId,
		NodePoolID: resp.Id,
	}, nil
}

func CreateBaseNodeLabels(params base_deploy_params.ClusterSetupRequest) map[string]string {
	labels := make(map[string]string)
	labels["org"] = fmt.Sprintf("%d", params.Ou.OrgID)

	if params.Cluster.ClusterName != "" {
		labels["app"] = params.Cluster.ClusterName
	}
	return labels
}

func (c *CreateSetupTopologyActivities) MakeNodePoolRequest(ctx context.Context, params base_deploy_params.ClusterSetupRequest) (do_types.DigitalOceanNodePoolRequestStatus, error) {
	taint := godo.Taint{
		Key:    fmt.Sprintf("org-%d", params.Ou.OrgID),
		Value:  fmt.Sprintf("org-%d", params.Ou.OrgID),
		Effect: "NoSchedule",
	}
	appTaint := godo.Taint{
		Key:    "app",
		Value:  params.Cluster.ClusterName,
		Effect: "NoSchedule",
	}
	labels := CreateBaseNodeLabels(params)
	tmp := strings.Split(params.CloudCtxNs.Namespace, "-")
	suffix := tmp[len(tmp)-1]
	taints := []godo.Taint{taint}
	if params.AppTaint {
		taints = append(taints, appTaint)
	}
	if strings.HasPrefix(params.Nodes.Slug, "so") {
		labels = hestia_digitalocean.AddDoNvmeLabels(labels)
	}
	nodePoolName := strings.ToLower(fmt.Sprintf("do-%d-%s", params.Ou.OrgID, suffix))
	if len(nodePoolName) > 39 {
		nodePoolName = nodePoolName[:39]
	}
	log.Info().Interface("nodePoolName", nodePoolName).Msg("MakeNodePoolRequest")
	nodesReq := &godo.KubernetesNodePoolCreateRequest{
		Name:   nodePoolName,
		Size:   params.Nodes.Slug,
		Count:  int(params.NodesQuantity),
		Labels: labels,
		Taints: taints,
	}
	// TODO remove hard code cluster id
	clusterID := "0de1ee8e-7b90-45ea-b966-e2d2b7976cf9"
	node, err := api_auth_temporal.DigitalOcean.CreateNodePool(ctx, clusterID, nodesReq)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Info().Interface("nodesReq", nodesReq).Msg("CreateNodePool already exists")
			return do_types.DigitalOceanNodePoolRequestStatus{
				ClusterID:  clusterID,
				NodePoolID: node.ID,
			}, nil
		}
		log.Err(err).Interface("nodes", params.Nodes).Msg("CreateNodePool error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	return do_types.DigitalOceanNodePoolRequestStatus{
		ClusterID:  clusterID,
		NodePoolID: node.ID,
	}, nil
}

func (c *CreateSetupTopologyActivities) SelectOvhNodeResources(ctx context.Context, request base_deploy_params.DestroyResourcesRequest) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	log.Info().Interface("request", request).Msg("SelectOvhNodeResources")
	nps, err := hestia_compute_resources.OvhSelectNodeResources(ctx, request.Ou.OrgID, request.OrgResourceIDs)
	if err != nil {
		log.Err(err).Interface("request", request).Msg("SelectOvhNodeResources: OvhSelectNodeResources error")
		return nps, err
	}
	return nps, err
}

func (c *CreateSetupTopologyActivities) SelectNodeResources(ctx context.Context, request base_deploy_params.DestroyResourcesRequest) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	log.Info().Interface("request", request).Msg("SelectNodeResources")
	nps, err := hestia_compute_resources.SelectNodeResources(ctx, request.Ou.OrgID, request.OrgResourceIDs)
	if err != nil {
		log.Err(err).Interface("request", request).Msg("SelectNodeResources: SelectNodeResources error")
		return nps, err
	}
	return nps, err
}

func (c *CreateSetupTopologyActivities) EndResourceService(ctx context.Context, request base_deploy_params.DestroyResourcesRequest) error {
	log.Info().Interface("request", request).Msg("EndResourceService")
	err := hestia_compute_resources.UpdateEndServiceOrgResources(ctx, request.Ou.OrgID, request.OrgResourceIDs)
	if err != nil {
		log.Err(err).Interface("request", request).Msg("EndResourceService: UpdateEndServiceOrgResources error")
		return err
	}
	return err
}

// clusterID := "0de1ee8e-7b90-45ea-b966-e2d2b7976cf9"

func (c *CreateSetupTopologyActivities) RemoveNodePoolRequest(ctx context.Context, nodePool do_types.DigitalOceanNodePoolRequestStatus) error {
	log.Info().Interface("nodePool", nodePool).Msg("RemoveNodePoolRequest")
	err := api_auth_temporal.DigitalOcean.RemoveNodePool(ctx, nodePool.ClusterID, nodePool.NodePoolID)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			log.Info().Interface("nodePool", nodePool).Msg("RemoveNodePoolRequest: node pool not found")
			return nil
		}
		log.Err(err).Interface("nodePool", nodePool).Msg("RemoveNodePool error")
		return err
	}
	return err
}

func (c *CreateSetupTopologyActivities) OvhRemoveNodePoolRequest(ctx context.Context, nodePool do_types.DigitalOceanNodePoolRequestStatus) error {
	err := api_auth_temporal.OvhCloud.RemoveNodePool(ctx, nodePool.ClusterID, nodePool.NodePoolID)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			log.Info().Interface("nodePool", nodePool).Msg("OvhRemoveNodePoolRequest: node pool not found")
			return nil
		}
		log.Err(err).Interface("nodePool", nodePool).Msg("OvhRemoveNodePoolRequest error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) RemoveFreeTrialOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest) error {
	err := hestia_compute_resources.RemoveFreeTrialOrgResources(ctx, params.Ou.OrgID)
	if err != nil {
		log.Err(err).Interface("ou", params.Ou).Msg("RemoveFreeTrialOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) UpdateFreeTrialOrgResourcesToPaid(ctx context.Context, params base_deploy_params.ClusterSetupRequest) error {
	err := hestia_compute_resources.UpdateFreeTrialOrgResourcesToPaid(ctx, params.Ou.OrgID)
	if err != nil {
		log.Err(err).Interface("ou", params.Ou).Msg("RemoveFreeTrialOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) SelectFreeTrialNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	nps, err := hestia_compute_resources.SelectFreeTrialDigitalOceanNodes(ctx, orgID)
	if err != nil {
		log.Err(err).Int("orgID", orgID).Msg("SelectFreeTrialNodes: SelectFreeTrialDigitalOceanNodes error")
		return nps, err
	}
	return nps, err
}

func (c *CreateSetupTopologyActivities) OvhSelectFreeTrialNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	eksNps, err := hestia_compute_resources.OvhSelectFreeTrialNodes(ctx, orgID)
	if err != nil {
		log.Err(err).Int("orgID", orgID).Msg("OvhSelectFreeTrialNodes: OvhSelectFreeTrialNodes error")
		return eksNps, err
	}
	return eksNps, err
}
