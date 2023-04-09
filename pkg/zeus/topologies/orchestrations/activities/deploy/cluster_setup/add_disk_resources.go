package deploy_topology_activities_create_setup

import (
	"context"

	"github.com/rs/zerolog/log"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
)

func (c *CreateSetupTopologyActivities) AddDiskResourcesToOrg(ctx context.Context, orgID int, disks hestia_autogen_bases.Disks, quantity float64) error {
	err := hestia_compute_resources.AddResourcesToOrg(ctx, orgID, disks.ResourceID, quantity)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("disks", disks).Msg("AddDiskResourcesToOrg error")
		return err
	}
	return nil
}
