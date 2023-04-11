package deploy_topology_activities_create_setup

import (
	"context"

	"github.com/rs/zerolog/log"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"k8s.io/apimachinery/pkg/api/resource"
)

func (c *CreateSetupTopologyActivities) AddDiskResourcesToOrg(ctx context.Context, params base_deploy_params.ClusterSetupRequest, disk hestia_autogen_bases.Disks) error {
	q, err := digitalOceanBlockStorageBillingUnits(ctx, disk.DiskUnits)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("disks", params.Disks).Msg("AddDiskResourcesToOrg error")
		return err
	}
	err = hestia_compute_resources.AddResourcesToOrg(ctx, params.Ou.OrgID, disk.ResourceID, q, params.FreeTrial)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("disks", params.Disks).Msg("AddDiskResourcesToOrg error")
		return err
	}
	return nil
}

func digitalOceanBlockStorageBillingUnits(ctx context.Context, qtyString string) (float64, error) {
	r, err := resource.ParseQuantity(qtyString)
	if err != nil {
		log.Error().Err(err).Msg("resource.ParseQuantity")
		return 0, err
	}
	rawValue := r.Value()
	q := resource.NewQuantity(rawValue, resource.BinarySI)
	q.ScaledValue(resource.Mega)
	miValue := q.AsDec().UnscaledBig().Int64() / (1024 * 1024 * 1024)
	billableUnits := float64(miValue / 100)
	return billableUnits, nil
}
