package deploy_topology_activities_create_setup

import (
	"context"

	"github.com/rs/zerolog/log"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

func (c *CreateSetupTopologyActivities) AddDomainRecord(ctx context.Context, cloudCtxNs zeus_common_types.CloudCtxNs) error {
	dr, err := api_auth_temporal.DigitalOcean.CreateDomain(ctx, cloudCtxNs)
	if err != nil {
		log.Ctx(ctx).Error().Interface("dr", dr).Err(err).Msg("failed to create domain record")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) RemoveDomainRecord(ctx context.Context, cloudCtxNs zeus_common_types.CloudCtxNs) error {
	err := api_auth_temporal.DigitalOcean.RemoveSubDomainARecord(ctx, cloudCtxNs)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to remove domain record")
		return err
	}
	return nil
}
