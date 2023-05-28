package deploy_topology_update_activities

import (
	"context"

	"github.com/rs/zerolog/log"
	read_topologies "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies"
	read_topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/definitions/state"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
)

type TopologyUpdateActivity struct {
	Host string
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *TopologyUpdateActivity) GetActivities() ActivitiesSlice {
	return []interface{}{d.GetClustersToUpdate, d.GetClusterTopologyAtCloudCtxNs}
}

func (d *TopologyUpdateActivity) GetClustersToUpdate(ctx context.Context, params base_deploy_params.FleetUpgradeWorkflowRequest) ([]read_topologies.ClusterAppView, error) {
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
