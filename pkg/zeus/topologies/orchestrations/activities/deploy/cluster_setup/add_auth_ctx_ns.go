package deploy_topology_activities_create_setup

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	create_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology"
)

func (c *CreateSetupTopologyActivities) AddAuthCtxNsOrg(ctx context.Context, newCloudCtxAuth create_topology.CreateTopologiesOrgCloudCtxNs) error {
	uuidNamespace, err := uuid.NewUUID()
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("cloudAuth", newCloudCtxAuth).Msg("AddAuthCtxNsOrg: NewUUID error")
		return err
	}
	newCloudCtxAuth.Namespace = uuidNamespace.String()
	err = newCloudCtxAuth.InsertTopologyAccessCloudCtxNs(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("cloudAuth", newCloudCtxAuth).Msg("AddAuthCtxNsOrg: InsertTopologyAccessCloudCtxNs error")
		return err
	}
	return err
}
