package deploy_topology_activities_create_setup

import (
	"context"

	"github.com/rs/zerolog/log"
	create_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

func (c *CreateSetupTopologyActivities) AddAuthCtxNsOrg(ctx context.Context, params base_deploy_params.ClusterSetupRequest) error {
	newCloudCtxAuth := create_topology.CreateTopologiesOrgCloudCtxNs{}
	err := newCloudCtxAuth.InsertTopologyAccessCloudCtxNs(ctx, params.Ou.OrgID, params.CloudCtxNs)
	if err != nil {
		log.Err(err).Interface("cloudAuth", newCloudCtxAuth).Msg("AddAuthCtxNsOrg: InsertTopologyAccessCloudCtxNs error")
		return err
	}
	return err
}

func (c *CreateSetupTopologyActivities) RemoveAuthCtxNsOrg(ctx context.Context, orgID int, cloudCtxNs zeus_common_types.CloudCtxNs) error {
	err := create_topology.DeleteTopologyAccessCloudCtxNs(ctx, orgID, cloudCtxNs)
	if err != nil {
		log.Err(err).Interface("cloudCtxNs", cloudCtxNs).Msg("RemoveAuthCtxNsOrg: DeleteTopologyAccessCloudCtxNs error")
		return err
	}
	return err
}
