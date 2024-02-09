package deploy_topology_activities_create_setup

import (
	"context"
	"fmt"
	ht "net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/smithy-go"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	hestia_eks_aws "github.com/zeus-fyi/olympus/pkg/hestia/aws"
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"k8s.io/apimachinery/pkg/util/validation"
)

func (c *CreateSetupTopologyActivities) EksAddNodePoolToOrgResources(ctx context.Context, params base_deploy_params.ClusterSetupRequest, npStatus do_types.DigitalOceanNodePoolRequestStatus) error {
	err := hestia_compute_resources.AddEksNodePoolResourcesToOrg(ctx, params.Ou.OrgID, params.Nodes.ResourceID, params.NodesQuantity, npStatus.NodePoolID, npStatus.ClusterID, params.FreeTrial, params.CloudCtxNs.ClusterCfgStrID)
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

func (c *CreateSetupTopologyActivities) PrivateEksMakeNodePoolRequest(ctx context.Context, params base_deploy_params.ClusterSetupRequest) (do_types.DigitalOceanNodePoolRequestStatus, error) {
	cfgID, err := strconv.Atoi(params.CloudCtxNs.ClusterCfgStrID)
	if err != nil {
		log.Err(err).Interface("cloudCtxNs", params.CloudCtxNs).Msg("PrivateEksMakeNodePoolRequest: strconv.Atoi error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	if params.CloudCtxNs.Context == "" {
		log.Err(fmt.Errorf("context is empty")).Interface("ou", params.Ou).Msg("PrivateEksMakeNodePoolRequest: context is empty")
		return do_types.DigitalOceanNodePoolRequestStatus{}, fmt.Errorf("context is empty")
	}
	ps, err := aws_secrets.GetServiceAccountSecrets(ctx, params.Ou)
	if err != nil {
		log.Err(err).Interface("ou", params.Ou).Msg("PrivateEksMakeNodePoolRequest: GetServiceAccountSecrets error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	var clusterName string
	var eksServiceAuth aegis_aws_auth.AuthAWS
	for cn, v := range ps.AwsEksServiceMap {
		if v.Region == params.CloudCtxNs.Region {
			eksServiceAuth = v
			clusterName = cn
			break
		}
	}
	if eksServiceAuth.Region == "" || eksServiceAuth.AccessKey == "" || eksServiceAuth.SecretKey == "" || clusterName == "" {
		log.Err(err).Interface("ou", params.Ou).Msg("PrivateEksMakeNodePoolRequest: GetServiceAccountSecrets error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	eksCredsAuth := hestia_eks_aws.EksCredentials{
		Creds:       eksServiceAuth,
		ClusterName: clusterName,
	}
	kubeConfig, err := hestia_eks_aws.GetEksKubeConfig(ctx, eksCredsAuth)
	if err != nil {
		log.Err(err).Interface("ou", params.Ou).Msg("PrivateEksMakeNodePoolRequest: GetEksKubeConfig error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	subnets, err := kubeConfig.GetEksSubnets()
	if err != nil {
		log.Err(err).Interface("ou", params.Ou).Msg("PrivateEksMakeNodePoolRequest: GetEksSubnets error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	role, err := kubeConfig.GetEksRoleArn()
	if err != nil {
		log.Err(err).Interface("ou", params.Ou).Msg("PrivateEksMakeNodePoolRequest: GetEksRoleArn error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	labels := CreateBaseNodeLabels(params)
	//orgTaint := types.Taint{
	//	Effect: "NO_SCHEDULE",
	//	Key:    aws.String(fmt.Sprintf("org-%d", params.Ou.OrgID)),
	//	Value:  aws.String(fmt.Sprintf("org-%d", params.Ou.OrgID)),
	//}
	var taints []types.Taint
	if len(params.Cluster.ClusterName) > 0 && params.AppTaint {
		appTaint := types.Taint{
			Effect: "NO_SCHEDULE",
			Key:    aws.String("app"),
			Value:  aws.String(strings.ToLower(params.Cluster.ClusterName)),
		}
		taints = append(taints, appTaint)
	}
	nsAlias := NamespaceAlias(params.Cluster.ClusterName)
	nodeGroupName := strings.ToLower(fmt.Sprintf("aws-%d-%s-z", params.Ou.OrgID, nsAlias))
	if len(nodeGroupName) > 39 {
		nodeGroupName = nodeGroupName[:38] + "z"
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
		ClusterName:        aws.String(clusterName),
		NodeRole:           role,
		NodegroupName:      aws.String(nodeGroupName),
		AmiType:            types.AMITypesAl2X8664,
		Subnets:            subnets,
		ClientRequestToken: aws.String(nodeGroupName),
		InstanceTypes:      instanceTypes,
		LaunchTemplate:     lt,
		Labels:             labels,
		ScalingConfig: &types.NodegroupScalingConfig{
			DesiredSize: aws.Int32(int32(params.NodesQuantity)),
			MaxSize:     aws.Int32(int32(params.NodesQuantity)),
			MinSize:     aws.Int32(int32(params.NodesQuantity)),
		},
		Taints: taints,
		Tags:   labels,
	}

	eka, err := hestia_eks_aws.InitAwsEKS(ctx, eksServiceAuth)
	if err != nil {
		log.Err(err).Msg("GetKubeConfig: failed to init EKS client")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	_, err = eka.AddNodeGroup(ctx, nr)
	if err != nil {
		log.Err(err).Interface("nodes", params.Nodes).Msg("PrivateEksMakeNodePoolRequest error")
		return do_types.DigitalOceanNodePoolRequestStatus{}, err
	}
	return do_types.DigitalOceanNodePoolRequestStatus{
		ExtClusterCfgID: cfgID,
		ClusterID:       params.CloudCtxNs.Context,
		NodePoolID:      nodeGroupName,
	}, nil
}

func (c *CreateSetupTopologyActivities) EksMakeNodePoolRequest(ctx context.Context, params base_deploy_params.ClusterSetupRequest) (do_types.DigitalOceanNodePoolRequestStatus, error) {
	labels := CreateBaseNodeLabels(params)
	orgTaint := types.Taint{
		Effect: "NO_SCHEDULE",
		Key:    aws.String(fmt.Sprintf("org-%d", params.Ou.OrgID)),
		Value:  aws.String(fmt.Sprintf("org-%d", params.Ou.OrgID)),
	}
	taints := []types.Taint{
		orgTaint,
	}
	if len(params.Cluster.ClusterName) > 0 && params.AppTaint {
		appTaint := types.Taint{
			Effect: "NO_SCHEDULE",
			Key:    aws.String("app"),
			Value:  aws.String(strings.ToLower(params.Cluster.ClusterName)),
		}
		taints = append(taints, appTaint)
	}

	nsAlias := NamespaceAlias(params.Cluster.ClusterName)
	nodeGroupName := strings.ToLower(fmt.Sprintf("aws-%d-%s-z", params.Ou.OrgID, nsAlias))
	if len(nodeGroupName) > 39 {
		nodeGroupName = nodeGroupName[:38] + "z"
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

		//errSmithy, ok := err.(*smithy.OperationError)
		//if !ok {
		//	log.Err(err).Interface("nodes", params.Nodes).Msg("EksMakeNodePoolRequest error")
		//	return do_types.DigitalOceanNodePoolRequestStatus{}, err
		//}
		//httpErr := errSmithy.Err.(*http.ResponseError)
		//httpResponse := httpErr.HTTPStatusCode()
		//if httpResponse == ht.StatusConflict {
		//	log.Info().Interface("nodeGroup", nodeGroupName).Msg("EksMakeNodePoolRequest already exists")
		//	return do_types.DigitalOceanNodePoolRequestStatus{
		//		ClusterID:  hestia_eks_aws.AwsUsWest1Context,
		//		NodePoolID: nodeGroupName,
		//	}, nil
		//}
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
		log.Err(err).Interface("nodePool", nodePool).Msg("EksRemoveNodePoolRequest error")
		errSmithy, ok := err.(*smithy.OperationError)
		if ok {
			httpErr, ok2 := errSmithy.Err.(*http.ResponseError)
			if ok2 {
				httpResponse := httpErr.HTTPStatusCode()
				if httpResponse == ht.StatusConflict || httpResponse == ht.StatusNotFound {
					log.Info().Interface("nodePool", nodePool).Msg("EksRemoveNodePoolRequest: node pool not found")
					return nil
				}
			}
		}
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) PrivateEksRemoveNodePoolRequest(ctx context.Context, ou org_users.OrgUser, cloudCtxNs zeus_common_types.CloudCtxNs, nodePool do_types.DigitalOceanNodePoolRequestStatus) error {
	nr := &eks.DeleteNodegroupInput{
		ClusterName:   aws.String(nodePool.ClusterID),
		NodegroupName: aws.String(nodePool.NodePoolID),
	}
	ps, err := aws_secrets.GetServiceAccountSecrets(ctx, ou)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("PrivateEksRemoveNodePoolRequest: GetServiceAccountSecrets error")
		return err
	}
	var clusterName string
	var eksServiceAuth aegis_aws_auth.AuthAWS
	for cn, v := range ps.AwsEksServiceMap {
		if v.Region == cloudCtxNs.Region && cn == nodePool.ClusterID {
			eksServiceAuth = v
			clusterName = cn
			break
		}
	}
	if eksServiceAuth.Region == "" || eksServiceAuth.AccessKey == "" || eksServiceAuth.SecretKey == "" || clusterName == "" {
		log.Err(err).Interface("ou", ou).Msg("PrivateEksRemoveNodePoolRequest: GetServiceAccountSecrets error")
		return err
	}
	if clusterName != nodePool.ClusterID {
		log.Err(fmt.Errorf("clusterName and nodePool.ClusterID do not match")).Interface("ou", ou).Interface("nodePool", nodePool).Msg("PrivateEksRemoveNodePoolRequest: clusterName and nodePool.ClusterID do not match")
		return fmt.Errorf("clusterName and nodePool.ClusterID do not match")
	}
	eka, err := hestia_eks_aws.InitAwsEKS(ctx, eksServiceAuth)
	if err != nil {
		log.Err(err).Msg("GetKubeConfig: failed to init EKS client")
		return err
	}
	_, err = eka.RemoveNodeGroup(ctx, nr)
	if err != nil {
		log.Err(err).Interface("nodePool", nodePool).Msg("PrivateEksRemoveNodePoolRequest error")
		errSmithy, ok := err.(*smithy.OperationError)
		if ok {
			httpErr, ok2 := errSmithy.Err.(*http.ResponseError)
			if ok2 {
				httpResponse := httpErr.HTTPStatusCode()
				if httpResponse == ht.StatusConflict || httpResponse == ht.StatusNotFound {
					log.Info().Interface("nodePool", nodePool).Msg("PrivateEksRemoveNodePoolRequest: node pool not found")
					return nil
				}
			}
		}
		return err
	}
	return nil
}

func NamespaceAlias(clusterName string) string {
	clusterID := uuid.New()
	suffix := strings.Split(clusterID.String(), "-")[0]
	if len(clusterName) <= 0 {
		return "z" + suffix + "z"
	}
	alias := fmt.Sprintf("%s-%s", clusterName, suffix)
	clusterNs := clusterID.String()
	if validation.IsDNS1123Label(alias) == nil {
		clusterNs = alias
	} else {
		clusterNs = suffix
	}
	return clusterNs
}
