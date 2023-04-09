package deploy_topology_activities_create_setup

import (
	"context"

	"github.com/rs/zerolog/log"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	deploy_workflow_cluster_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create_setup"
)

func (c *CreateSetupTopologyActivities) AddDiskResourcesToOrg(ctx context.Context, params deploy_workflow_cluster_setup.ClusterSetupRequest) error {
	err := hestia_compute_resources.AddResourcesToOrg(ctx, params.Ou.OrgID, params.Disks.ResourceID, params.DisksQuantity)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("disks", params.Disks).Msg("AddDiskResourcesToOrg error")
		return err
	}
	return nil
}
