package quicknode_orchestrations

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	hestia_quicknode_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/quiknode"
	platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/iris/orchestrations"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode"
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
		h.Provision, h.UpdateProvision, h.Deprovision, h.Deactivate, h.DeprovisionCache,
	}
}

func (h *HestiaQuicknodeActivities) Provision(ctx context.Context, pr hestia_quicknode.ProvisionRequest, ou org_users.OrgUser, user hestia_quicknode.QuickNodeUserInfo) error {
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
		WssURL: sql.NullString{
			String: pr.WssUrl,
			Valid:  len(pr.WssUrl) > 0,
		},
		Chain: sql.NullString{
			String: pr.Chain,
			Valid:  len(pr.Chain) > 0,
		},
	}

	cas := make([]hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses, len(pr.ContractAddresses))
	for i, ca := range pr.ContractAddresses {
		cas[i] = hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses{
			ContractAddress: ca,
		}
	}
	car := make([]hestia_autogen_bases.ProvisionedQuicknodeServicesReferers, len(pr.ContractAddresses))
	for i, re := range pr.Referers {
		car[i] = hestia_autogen_bases.ProvisionedQuicknodeServicesReferers{
			Referer: re,
		}
	}
	qs := hestia_quicknode_models.QuickNodeService{
		IsTest:                       user.IsTest,
		ProvisionedQuickNodeServices: ps,
		ProvisionedQuicknodeServicesContractAddresses: cas,
		ProvisionedQuicknodeServicesReferers:          car,
	}

	err := hestia_quicknode_models.InsertProvisionedQuickNodeService(ctx, qs)
	if err != nil {
		log.Warn().Interface("ou", ou).Err(err).Msg("Provision: InsertProvisionedQuickNodeService")
		return err
	}
	if !user.Verified {
		co := create_org_users.NewCreateOrgUser()
		err = co.InsertOrgUserWithNewQuickNodeKeyForService(ctx, qs.QuickNodeID)
		if err != nil {
			log.Warn().Interface("ou", ou).Err(err).Msg("Provision: InsertOrgUserWithNewQuickNodeKeyForService")
			return err
		}
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
	cas := make([]hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses, len(pr.ContractAddresses))
	for i, ca := range pr.ContractAddresses {
		cas[i] = hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses{
			ContractAddress: ca,
		}
	}
	car := make([]hestia_autogen_bases.ProvisionedQuicknodeServicesReferers, len(pr.ContractAddresses))
	for i, re := range pr.Referers {
		car[i] = hestia_autogen_bases.ProvisionedQuicknodeServicesReferers{
			Referer: re,
		}
	}
	qs := hestia_quicknode_models.QuickNodeService{
		ProvisionedQuickNodeServices:                  ps,
		ProvisionedQuicknodeServicesContractAddresses: cas,
		ProvisionedQuicknodeServicesReferers:          car,
	}
	err := hestia_quicknode_models.UpdateProvisionedQuickNodeService(ctx, qs)
	if err != nil {
		log.Warn().Interface("ou", ou).Err(err).Msg("Provision: UpdateProvision")
		return err
	}
	return nil
}

func (h *HestiaQuicknodeActivities) Deprovision(ctx context.Context, dp hestia_quicknode.DeprovisionRequest, ou org_users.OrgUser) error {
	err := hestia_quicknode_models.DeprovisionQuickNodeServices(ctx, dp.QuickNodeID)
	if err != nil {
		log.Warn().Interface("ou", ou).Err(err).Msg("Provision: Deprovision")
		return err
	}
	return nil
}

func (h *HestiaQuicknodeActivities) DeprovisionCache(ctx context.Context, ou org_users.OrgUser) error {
	pr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou: ou,
	}
	err := platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisRemoveAllOrgRoutesFromCacheWorkflow(ctx, pr)
	if err != nil {
		log.Warn().Interface("ou", ou).Err(err).Msg("Provision: DeprovisionCache")
		return err
	}
	return nil
}

func (h *HestiaQuicknodeActivities) Deactivate(ctx context.Context, da hestia_quicknode.DeactivateRequest, ou org_users.OrgUser) error {
	err := hestia_quicknode_models.DeactivateProvisionedQuickNodeServiceEndpoint(ctx, da.QuickNodeID, da.EndpointID)
	if err != nil {
		log.Warn().Interface("ou", ou).Err(err).Msg("Provision: Deactivate")
		return err
	}
	return nil
}
