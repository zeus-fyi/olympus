package deploy_topology_activities_create_setup

import (
	"context"
	"fmt"
	ht "net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/smithy-go"
	"github.com/digitalocean/godo"
	"github.com/rs/zerolog/log"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	hestia_eks_aws "github.com/zeus-fyi/olympus/pkg/hestia/aws"
	hestia_digitalocean "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean"
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	hestia_gcp "github.com/zeus-fyi/olympus/pkg/hestia/gcp"
	hestia_ovhcloud "github.com/zeus-fyi/olympus/pkg/hestia/ovhcloud"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"google.golang.org/api/container/v1"
)

func (c *CreateSetupTopologyActivities) AddNodePoolToOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest, npStatus do_types.DigitalOceanNodePoolRequestStatus) error {
	err := hestia_compute_resources.AddDigitalOceanNodePoolResourcesToOrg(ctx, params.Ou.OrgID, params.Nodes.ResourceID, params.NodesQuantity, npStatus.NodePoolID, npStatus.ClusterID, params.FreeTrial)
	if err != nil {
		log.Err(err).Interface("nodes", params.Nodes).Msg("AddNodePoolToOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) GkeAddNodePoolToOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest, npStatus do_types.DigitalOceanNodePoolRequestStatus) error {
	err := hestia_compute_resources.AddGkeNodePoolResourcesToOrg(ctx, params.Ou.OrgID, params.Nodes.ResourceID, params.NodesQuantity, npStatus.NodePoolID, npStatus.ClusterID, params.FreeTrial)
	if err != nil {
		log.Err(err).Interface("nodes", params.Nodes).Msg("GkeAddNodePoolToOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) EksAddNodePoolToOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest, npStatus do_types.DigitalOceanNodePoolRequestStatus) error {
	err := hestia_compute_resources.AddEksNodePoolResourcesToOrg(ctx, params.Ou.OrgID, params.Nodes.ResourceID, params.NodesQuantity, npStatus.NodePoolID, npStatus.ClusterID, params.FreeTrial)
	if err != nil {
		log.Err(err).Interface("nodes", params.Nodes).Msg("EksAddNodePoolToOrgResources error")
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
	tmp := strings.Split(params.Namespace, "-")
	suffix := tmp[len(tmp)-1]
	nodeGroupName := fmt.Sprintf("nodepool-%d-%s", params.Ou.OrgID, suffix)
	label := make(map[string]string)
	label["org"] = fmt.Sprintf("%d", params.Ou.OrgID)
	label["app"] = params.Cluster.ClusterName
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
					Labels:      label,
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

func (c *CreateSetupTopologyActivities) EksMakeNodePoolRequest(ctx context.Context, params base_deploy_params.ClusterSetupRequest) (do_types.DigitalOceanNodePoolRequestStatus, error) {
	labels := make(map[string]string)
	labels["org"] = fmt.Sprintf("%d", params.Ou.OrgID)
	labels["app"] = params.Cluster.ClusterName

	tmp := strings.Split(params.Namespace, "-")
	suffix := tmp[len(tmp)-1]
	orgTaint := types.Taint{
		Effect: "NO_SCHEDULE",
		Key:    aws.String(fmt.Sprintf("org-%d", params.Ou.OrgID)),
		Value:  aws.String(fmt.Sprintf("org-%d", params.Ou.OrgID)),
	}
	appTaint := types.Taint{
		Effect: "NO_SCHEDULE",
		Key:    aws.String("app"),
		Value:  aws.String(params.Cluster.ClusterName),
	}
	taints := []types.Taint{
		orgTaint,
	}
	if params.AppTaint {
		taints = append(taints, appTaint)
	}

	var lt *types.LaunchTemplateSpecification
	instanceTypes := []string{params.Nodes.Slug}
	nodeGroupName := fmt.Sprintf("nodepool-%d-%s", params.Ou.OrgID, suffix)
	id := hestia_eks_aws.SlugToLaunchTemplateID[params.Nodes.Slug]
	if id != "" {
		labels = hestia_eks_aws.AddAwsEksNvmeLabels(labels)
		instanceTypes = nil
		st := hestia_eks_aws.SlugToInstanceTemplateName[params.Nodes.Slug]
		lt = hestia_eks_aws.GetLaunchTemplate(id, st)
	}

	nr := &eks.CreateNodegroupInput{
		ClusterName:        aws.String(hestia_eks_aws.AwsUsWest1Context),
		NodeRole:           aws.String(hestia_eks_aws.AwsEksRole),
		NodegroupName:      aws.String(nodeGroupName),
		AmiType:            types.AMITypesAl2X8664,
		Subnets:            hestia_eks_aws.UsWestSubnetIDs,
		ClientRequestToken: aws.String(nodeGroupName),
		InstanceTypes:      instanceTypes,
		LaunchTemplate:     lt,
		Labels:             labels,
		ReleaseVersion:     nil,
		ScalingConfig: &types.NodegroupScalingConfig{
			DesiredSize: aws.Int32(int32(params.NodesQuantity)),
			MaxSize:     aws.Int32(int32(params.NodesQuantity)),
			MinSize:     aws.Int32(int32(params.NodesQuantity)),
		},
		Taints: taints,
	}
	_, err := api_auth_temporal.Eks.AddNodeGroup(ctx, nr)
	if err != nil {
		errSmithy := err.(*smithy.OperationError)
		httpErr := errSmithy.Err.(*http.ResponseError)
		httpResponse := httpErr.HTTPStatusCode()
		if httpResponse == ht.StatusConflict {
			log.Info().Interface("nodeGroup", nodeGroupName).Msg("EksMakeNodePoolRequest already exists")
			return do_types.DigitalOceanNodePoolRequestStatus{
				ClusterID:  hestia_eks_aws.AwsUsWest1Context,
				NodePoolID: nodeGroupName,
			}, nil
		}
		log.Err(err).Interface("nodes", params.Nodes).Msg("EksMakeNodePoolRequest error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	return do_types.DigitalOceanNodePoolRequestStatus{
		ClusterID:  hestia_eks_aws.AwsUsWest1Context,
		NodePoolID: nodeGroupName,
	}, nil
}

func (c *CreateSetupTopologyActivities) GkeMakeNodePoolRequest(ctx context.Context, params base_deploy_params.ClusterSetupRequest) (do_types.DigitalOceanNodePoolRequestStatus, error) {
	labels := make(map[string]string)
	labels["org"] = fmt.Sprintf("%d", params.Ou.OrgID)
	labels["app"] = params.Cluster.ClusterName

	tmp := strings.Split(params.Namespace, "-")
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
	ni := hestia_gcp.GkeNodePoolInfo{
		Name:             fmt.Sprintf("nodepool-%d-%s", params.Ou.OrgID, suffix),
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
		log.Err(err).Interface("nodes", params.Nodes).Msg("GkeMakeNodePoolRequest error")
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
	appTaint := godo.Taint{
		Key:    "app",
		Value:  params.Cluster.ClusterName,
		Effect: "NoSchedule",
	}
	labels := make(map[string]string)
	labels["org"] = fmt.Sprintf("%d", params.Ou.OrgID)
	labels["app"] = params.Cluster.ClusterName
	tmp := strings.Split(params.Namespace, "-")
	suffix := tmp[len(tmp)-1]
	taints := []godo.Taint{taint}
	if params.AppTaint {
		taints = append(taints, appTaint)
	}
	if strings.HasPrefix(params.Nodes.Slug, "so") {
		labels = hestia_digitalocean.AddDoNvmeLabels(labels)
	}
	nodePoolName := fmt.Sprintf("nodepool-%d-%s", params.Ou.OrgID, suffix)
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

func (c *CreateSetupTopologyActivities) SelectEksNodeResources(ctx context.Context, request base_deploy_params.DestroyResourcesRequest) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	log.Info().Interface("request", request).Msg("SelectEksNodeResources")
	nps, err := hestia_compute_resources.EksSelectNodeResources(ctx, request.Ou.OrgID, request.OrgResourceIDs)
	if err != nil {
		log.Err(err).Interface("request", request).Msg("SelectEksNodeResources: EksSelectNodeResources error")
		return nps, err
	}
	return nps, err
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

func (c *CreateSetupTopologyActivities) GkeRemoveNodePoolRequest(ctx context.Context, nodePool do_types.DigitalOceanNodePoolRequestStatus) error {
	log.Info().Interface("nodePool", nodePool).Msg("RemoveNodePoolRequest")
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
		if strings.Contains(err.Error(), "Not found") {
			log.Info().Interface("nodePool", nodePool).Msg("GkeRemoveNodePoolRequest: node pool not found")
			return nil
		}
		log.Err(err).Interface("nodePool", nodePool).Msg("GkeRemoveNodePoolRequest error")
		return err
	}
	return err
}

func (c *CreateSetupTopologyActivities) EksRemoveNodePoolRequest(ctx context.Context, nodePool do_types.DigitalOceanNodePoolRequestStatus) error {
	nr := &eks.DeleteNodegroupInput{
		ClusterName:   aws.String(nodePool.ClusterID),
		NodegroupName: aws.String(nodePool.NodePoolID),
	}
	_, err := api_auth_temporal.Eks.RemoveNodeGroup(ctx, nr)
	if err != nil {
		errSmithy := err.(*smithy.OperationError)
		httpErr := errSmithy.Err.(*http.ResponseError)
		httpResponse := httpErr.HTTPStatusCode()
		if httpResponse == ht.StatusConflict || httpResponse == ht.StatusNotFound {
			log.Info().Interface("nodePool", nodePool).Msg("EksRemoveNodePoolRequest: node pool not found")
			return nil
		} else {
			log.Err(err).Interface("nodePool", nodePool).Msg("EksRemoveNodePoolRequest error")
			return err
		}
	}
	return nil
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

func (c *CreateSetupTopologyActivities) GkeSelectFreeTrialNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	gkeNps, err := hestia_compute_resources.GkeSelectFreeTrialNodes(ctx, orgID)
	if err != nil {
		log.Err(err).Int("orgID", orgID).Msg("GkeSelectFreeTrialNodes: GkeSelectFreeTrialNodes error")
		return gkeNps, err
	}
	return gkeNps, err
}

func (c *CreateSetupTopologyActivities) EksSelectFreeTrialNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	eksNps, err := hestia_compute_resources.EksSelectFreeTrialNodes(ctx, orgID)
	if err != nil {
		log.Err(err).Int("orgID", orgID).Msg("EksSelectFreeTrialNodes: EksSelectFreeTrialNodes error")
		return eksNps, err
	}
	return eksNps, err
}

func (c *CreateSetupTopologyActivities) OvhSelectFreeTrialNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	eksNps, err := hestia_compute_resources.OvhSelectFreeTrialNodes(ctx, orgID)
	if err != nil {
		log.Err(err).Int("orgID", orgID).Msg("OvhSelectFreeTrialNodes: OvhSelectFreeTrialNodes error")
		return eksNps, err
	}
	return eksNps, err
}
