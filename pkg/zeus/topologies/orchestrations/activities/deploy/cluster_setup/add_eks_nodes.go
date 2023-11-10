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
	"github.com/rs/zerolog/log"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	hestia_eks_aws "github.com/zeus-fyi/olympus/pkg/hestia/aws"
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
)

func (c *CreateSetupTopologyActivities) EksAddNodePoolToOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest, npStatus do_types.DigitalOceanNodePoolRequestStatus) error {
	err := hestia_compute_resources.AddEksNodePoolResourcesToOrg(ctx, params.Ou.OrgID, params.Nodes.ResourceID, params.NodesQuantity, npStatus.NodePoolID, npStatus.ClusterID, params.FreeTrial)
	if err != nil {
		log.Err(err).Interface("nodes", params.Nodes).Msg("EksAddNodePoolToOrgResources error")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) EksSelectFreeTrialNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	eksNps, err := hestia_compute_resources.EksSelectFreeTrialNodes(ctx, orgID)
	if err != nil {
		log.Err(err).Int("orgID", orgID).Msg("EksSelectFreeTrialNodes: EksSelectFreeTrialNodes error")
		return eksNps, err
	}
	return eksNps, err
}

func (c *CreateSetupTopologyActivities) EksMakeNodePoolRequest(ctx context.Context, params base_deploy_params.ClusterSetupRequest) (do_types.DigitalOceanNodePoolRequestStatus, error) {
	labels := CreateBaseNodeLabels(params)
	tmp := strings.Split(params.CloudCtxNs.Namespace, "-")
	suffix := tmp[len(tmp)-1]
	orgTaint := types.Taint{
		Effect: "NO_SCHEDULE",
		Key:    aws.String(fmt.Sprintf("org-%d", params.Ou.OrgID)),
		Value:  aws.String(fmt.Sprintf("org-%d", params.Ou.OrgID)),
	}
	appTaint := types.Taint{
		Effect: "NO_SCHEDULE",
		Key:    aws.String("app"),
		Value:  aws.String(strings.ToLower(params.Cluster.ClusterName)),
	}
	taints := []types.Taint{
		orgTaint,
	}
	if params.AppTaint {
		taints = append(taints, appTaint)
	}

	nodeGroupName := strings.ToLower(fmt.Sprintf("aws-%d-%s", params.Ou.OrgID, suffix))
	if len(nodeGroupName) > 39 {
		nodeGroupName = nodeGroupName[:39]
	}
	var lt *types.LaunchTemplateSpecification
	instanceTypes := []string{params.Nodes.Slug}
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
		Tags:   labels,
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

func (c *CreateSetupTopologyActivities) SelectEksNodeResources(ctx context.Context, request base_deploy_params.DestroyResourcesRequest) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	log.Info().Interface("request", request).Msg("SelectEksNodeResources")
	nps, err := hestia_compute_resources.EksSelectNodeResources(ctx, request.Ou.OrgID, request.OrgResourceIDs)
	if err != nil {
		log.Err(err).Interface("request", request).Msg("SelectEksNodeResources: EksSelectNodeResources error")
		return nps, err
	}
	return nps, err
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
