package deploy_topology_activities_create_setup

import (
	"context"

	"github.com/rs/zerolog/log"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
)

func (c *CreateSetupTopologyActivities) AddDomainRecord(ctx context.Context, name string) error {
	dr, err := api_auth_temporal.DigitalOcean.CreateDomain(ctx, name)
	if err != nil {
		log.Ctx(ctx).Error().Interface("dr", dr).Err(err).Msg("failed to create domain record")
		return err
	}
	return nil
}

func (c *CreateSetupTopologyActivities) RemoveDomainRecord(ctx context.Context, name string) error {
	err := api_auth_temporal.DigitalOcean.RemoveSubDomainARecord(ctx, name)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to remove domain record")
		return err
	}
	return nil
}
