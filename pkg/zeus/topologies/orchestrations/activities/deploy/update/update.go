package deploy_topology_update_activities

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"
	read_topologies "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies"
	read_topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/definitions/state"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	topology_auths "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type TopologyUpdateActivity struct {
	Host string
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *TopologyUpdateActivity) GetActivities() ActivitiesSlice {
	return []interface{}{d.GetClustersToUpdate, d.GetClusterTopologyAtCloudCtxNs, d.DiffClusterUpdate, d.GetClustersToRolloutRestart,
		d.RestartWorkload}
}

const (
	deployment  = "deployment"
	statefulset = "statefulset"
)

func (d *TopologyUpdateActivity) RestartWorkload(ctx context.Context, cloudCtxNs zeus_common_types.CloudCtxNs, workloadName, workloadType string) error {
	if cloudCtxNs.CheckIfEmpty() {
		log.Warn().Interface("cloudCtxNs", cloudCtxNs).Msg("RestartWorkload: cloudCtxNs is empty")
		return nil
	}
	wt := strings.ToLower(workloadType)
	switch wt {
	case deployment:
		_, err := topology_auths.K8Util.RolloutRestartDeployment(ctx, cloudCtxNs, workloadName, nil)
		if err != nil {
			log.Err(err).Interface("cloudCtxNs", cloudCtxNs).Msg("RolloutRestartDeployment: error")
			return err
		}
	case statefulset:
		err := topology_auths.K8Util.RolloutRestartStatefulSet(ctx, cloudCtxNs, workloadName, nil)
		if err != nil {
			log.Err(err).Interface("cloudCtxNs", cloudCtxNs).Msg("RolloutRestartStatefulSets: error")
			return err
		}
	default:
		return nil
	}
	return nil
}

func (d *TopologyUpdateActivity) GetClustersToUpdate(ctx context.Context, params base_deploy_params.FleetUpgradeWorkflowRequest) ([]read_topologies.ClusterAppView, error) {
	resp, err := read_topologies.SelectClusterSingleAppView(ctx, params.OrgUser.OrgID, params.ClusterName)
	if err != nil {
		log.Err(err).Msg("GetClustersToUpdate")
		return resp, err
	}
	return resp, nil
}

func (d *TopologyUpdateActivity) GetClustersToRolloutRestart(ctx context.Context, params base_deploy_params.FleetRolloutRestartWorkflowRequest) ([]read_topologies.ClusterAppView, error) {
	resp, err := read_topologies.SelectClusterSingleAppView(ctx, params.OrgUser.OrgID, params.ClusterName)
	if err != nil {
		log.Err(err).Msg("GetClustersToUpdate")
		return resp, err
	}
	return resp, nil
}

func (d *TopologyUpdateActivity) GetClusterTopologyAtCloudCtxNs(ctx context.Context, params base_deploy_params.FleetUpgradeWorkflowRequest, clusterInfo read_topologies.ClusterAppView) (read_topology_deployment_status.ReadDeploymentStatusesGroup, error) {
	status := read_topology_deployment_status.NewReadDeploymentStatusesGroup()
	err := status.ReadLatestDeployedClusterTopologies(ctx, clusterInfo.CloudCtxNsID, params.OrgUser)
	if err != nil {
		log.Err(err).Interface("orgUser", params.OrgUser).Msg("GetClusterTopologyAtCloudCtxNs: ReadDeployedClusterTopologies")
		return status, err
	}
	return status, nil
}

func (d *TopologyUpdateActivity) DiffClusterUpdate(ctx context.Context, params base_deploy_params.FleetUpgradeWorkflowRequest, clusterInfo read_topologies.ClusterAppView, topGroup read_topology_deployment_status.ReadDeploymentStatusesGroup) (base_deploy_params.ClusterTopologyWorkflowRequest, error) {
	existingTopologyIDs, sbOptions := getSkeletonBaseNamesByClusterClassName(params.ClusterName, topGroup)
	m := make(map[string]map[string]bool)
	for _, val := range topGroup.Slice {
		if _, ok := m[val.ComponentBaseName]; !ok {
			m[val.ComponentBaseName] = make(map[string]bool)
		}
		m[val.ComponentBaseName][val.SkeletonBaseName] = true
	}
	cl, err := read_topology.SelectClusterTopologyFiltered(ctx, params.OrgUser.OrgID, params.ClusterName, sbOptions, m)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("DiffClusterUpdate: SelectClusterTopology")
		return base_deploy_params.ClusterTopologyWorkflowRequest{}, err
	}
	latestClTops := make(map[int]read_topology.ClusterTopologies)
	for _, top := range cl.Topologies {
		latestClTops[top.TopologyID] = top
	}
	var newTopIDs []int
	for key, v := range latestClTops {
		if _, exists := existingTopologyIDs[key]; !exists {
			newTopIDs = append(newTopIDs, key)
			log.Info().Str("clusterName", params.ClusterName).Interface("replacing", v).Msg("DiffClusterUpdate: replacing")
		}
	}
	if len(newTopIDs) == 0 {
		return base_deploy_params.ClusterTopologyWorkflowRequest{}, nil
	}
	clDeploy := base_deploy_params.ClusterTopologyWorkflowRequest{
		ClusterClassName: params.ClusterName,
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			CloudProvider: clusterInfo.CloudProvider,
			Region:        clusterInfo.Region,
			Context:       clusterInfo.Context,
			Namespace:     clusterInfo.Namespace,
			Alias:         clusterInfo.NamespaceAlias,
		},
		TopologyIDs: newTopIDs,
		OrgUser:     params.OrgUser,
		AppTaint:    params.AppTaint,
	}
	return clDeploy, nil
}

func getSkeletonBaseNamesByClusterClassName(clusterClassName string, topGroup read_topology_deployment_status.ReadDeploymentStatusesGroup) (map[int]bool, []string) {
	var names []string
	m := make(map[int]bool)
	for _, cluster := range topGroup.Slice {
		if cluster.ClusterName == clusterClassName {
			names = append(names, cluster.SkeletonBaseName)
			m[cluster.TopologyID] = true
		}
	}
	return m, names
}

func (d *TopologyUpdateActivity) GetClusterWorkloadsToRestart(ctx context.Context, params base_deploy_params.FleetRolloutRestartWorkflowRequest, topGroup read_topology_deployment_status.ReadDeploymentStatusesGroup) (*read_topology.ClusterTopology, error) {
	_, sbOptions := getSkeletonBaseNamesByClusterClassName(params.ClusterName, topGroup)
	m := make(map[string]map[string]bool)
	for _, val := range topGroup.Slice {
		if _, ok := m[val.ComponentBaseName]; !ok {
			m[val.ComponentBaseName] = make(map[string]bool)
		}
		m[val.ComponentBaseName][val.SkeletonBaseName] = true
	}
	cl, err := read_topology.SelectClusterTopologyFiltered(ctx, params.OrgUser.OrgID, params.ClusterName, sbOptions, m)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("DiffClusterUpdate: SelectClusterTopology")
		return nil, err
	}
	return &cl, nil
}
