package quicknode_orchestrations

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_quicknode_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/quiknode"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/quiknode"
)

type HestiaQuicknodeActivities struct {
}

func NewHestiaQuicknodeActivities() HestiaQuicknodeActivities {
	return HestiaQuicknodeActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (h *HestiaQuicknodeActivities) GetActivities() ActivitiesSlice {
	return []interface{}{
		h.Provision, h.UpdateProvision, h.Deprovision, h.Deprovision,
	}
}

func (h *HestiaQuicknodeActivities) Provision(ctx context.Context, pr hestia_quicknode.ProvisionRequest, ou org_users.OrgUser) error {
	ps := hestia_autogen_bases.ProvisionedQuickNodeServices{
		QuickNodeID: pr.QuickNodeID,
		EndpointID:  pr.EndpointID,
		HttpURL: sql.NullString{
			String: pr.HttpUrl,
			Valid:  len(pr.HttpUrl) > 0,
		},
		Network: sql.NullString{},
		Plan:    pr.Plan,
		Active:  true,
		OrgID:   ou.OrgID,
		WssURL: sql.NullString{
			String: pr.WssUrl,
			Valid:  len(pr.WssUrl) > 0,
		},
		Chain: sql.NullString{
			String: pr.Chain,
			Valid:  len(pr.Chain) > 0,
		},
	}
	err := hestia_quicknode_models.InsertProvisionedQuickNodeService(ctx, ps)
	if err != nil {
		log.Warn().Interface("ou", ou).Err(err).Msg("Provision: InsertProvisionedQuickNodeService")
		return err
	}
	return nil
}

func (h *HestiaQuicknodeActivities) UpdateProvision(ctx context.Context, pr hestia_quicknode.ProvisionRequest, ou org_users.OrgUser) error {
	ps := hestia_autogen_bases.ProvisionedQuickNodeServices{
		QuickNodeID: pr.QuickNodeID,
		EndpointID:  pr.EndpointID,
		HttpURL: sql.NullString{
			String: pr.HttpUrl,
			Valid:  len(pr.HttpUrl) > 0,
		},
		Network: sql.NullString{},
		Plan:    pr.Plan,
		Active:  true,
		OrgID:   ou.OrgID,
		WssURL: sql.NullString{
			String: pr.WssUrl,
			Valid:  len(pr.WssUrl) > 0,
		},
		Chain: sql.NullString{
			String: pr.Chain,
			Valid:  len(pr.Chain) > 0,
		},
	}
	err := hestia_quicknode_models.UpdateProvisionedQuickNodeService(ctx, ps)
	if err != nil {
		log.Warn().Interface("ou", ou).Err(err).Msg("Provision: InsertProvisionedQuickNodeService")
		return err
	}
	return nil
}

func (h *HestiaQuicknodeActivities) Deprovision(dp hestia_quicknode.DeprovisionRequest, ou org_users.OrgUser) error {

	return nil
}

func (h *HestiaQuicknodeActivities) Deactivate(da hestia_quicknode.DeactivateRequest, ou org_users.OrgUser) error {

	return nil
}
