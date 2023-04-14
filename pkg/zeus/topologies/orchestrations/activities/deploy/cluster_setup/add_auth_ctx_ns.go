package deploy_topology_activities_create_setup

import (
	"context"

	"github.com/rs/zerolog/log"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	create_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
)

func (c *CreateSetupTopologyActivities) AddAuthCtxNsOrg(ctx context.Context, params base_deploy_params.ClusterSetupRequest) error {
	newCloudCtxAuth := create_topology.CreateTopologiesOrgCloudCtxNs{
		TopologiesOrgCloudCtxNs: autogen_bases.TopologiesOrgCloudCtxNs{
			OrgID:          params.Ou.OrgID,
			CloudProvider:  params.CloudProvider,
			Context:        params.Context,
			Region:         params.Region,
			Namespace:      params.ClusterID.String(),
			NamespaceAlias: params.Namespace,
		}}
	err := newCloudCtxAuth.InsertTopologyAccessCloudCtxNs(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("cloudAuth", newCloudCtxAuth).Msg("AddAuthCtxNsOrg: InsertTopologyAccessCloudCtxNs error")
		return err
	}
	return err
}

func (c *CreateSetupTopologyActivities) RemoveAuthCtxNsOrg(ctx context.Context, params base_deploy_params.ClusterSetupRequest) error {
	newCloudCtxAuth := create_topology.CreateTopologiesOrgCloudCtxNs{
		TopologiesOrgCloudCtxNs: autogen_bases.TopologiesOrgCloudCtxNs{
			OrgID:         params.Ou.OrgID,
			CloudProvider: params.CloudProvider,
			Context:       params.Context,
			Region:        params.Region,
			Namespace:     params.Namespace,
		}}
	err := newCloudCtxAuth.DeleteTopologyAccessCloudCtxNs(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("cloudAuth", newCloudCtxAuth).Msg("RemoveAuthCtxNsOrg: DeleteTopologyAccessCloudCtxNs error")
		return err
	}
	return err
}
