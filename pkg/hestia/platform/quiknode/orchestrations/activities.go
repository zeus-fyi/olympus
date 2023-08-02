package quicknode_orchestrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	hestia_quicknode_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/quiknode"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/iris/orchestrations"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

const (
	IrisApiUrl = "https://iris.zeus.fyi"
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
		h.Provision, h.UpdateProvision, h.Deprovision, h.Deactivate, h.DeprovisionCache, h.CheckPlanOverages,
		h.IrisPlatformDeleteGroupTableCacheRequest, h.DeactivateApiKey, h.DeleteOrgGroupRoutingTable, h.InsertQuickNodeApiKey,
	}
}

func (h *HestiaQuicknodeActivities) DeleteOrgGroupRoutingTable(ctx context.Context, ou org_users.OrgUser, groupName string) error {
	err := iris_models.DeleteOrgGroupAndRoutes(context.Background(), ou.OrgID, groupName)
	if err != nil {
		log.Err(err).Msg("DeleteOrgGroupRoutingTable: DeleteOrgGroupRoutingTable")
		return err
	}
	return nil
}

func (h *HestiaQuicknodeActivities) InsertQuickNodeApiKey(ctx context.Context, pr hestia_quicknode.ProvisionRequest) error {
	co := create_org_users.NewCreateOrgUser()
	err := co.InsertOrgUserWithNewQuickNodeKeyForService(ctx, pr.QuickNodeID)
	if err != nil {
		log.Warn().Str("pr.QuickNodeID", pr.QuickNodeID).Err(err).Msg("Provision: InsertQuickNodeApiKey")
		return err
	}
	return nil
}

func (h *HestiaQuicknodeActivities) IrisPlatformDeleteGroupTableCacheRequest(ctx context.Context, ou org_users.OrgUser, groupName string) error {
	rc := resty_base.GetBaseRestyClient(IrisApiUrl, artemis_orchestration_auth.Bearer)
	refreshEndpoint := fmt.Sprintf("/v1/internal/router/delete/%d/%s", ou.OrgID, groupName)
	resp, err := rc.R().Get(refreshEndpoint)
	if err != nil {
		log.Err(err).Msg("HestiaQuicknodeActivities: IrisPlatformDeleteGroupTableCacheRequest")
		return err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Interface("orgUser", ou).Msg("HestiaQuicknodeActivities: IrisPlatformDeleteGroupTableCacheRequest")
		return err
	}
	return nil
}

func (h *HestiaQuicknodeActivities) Provision(ctx context.Context, ou org_users.OrgUser, pr hestia_quicknode.ProvisionRequest, user hestia_quicknode.QuickNodeUserInfo) error {
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
	if pr.Referers == nil {
		pr.Referers = []string{}
	}
	if pr.ContractAddresses == nil {
		pr.ContractAddresses = []string{}
	}
	cas := make([]hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses, len(pr.ContractAddresses))
	for i, ca := range pr.ContractAddresses {
		cas[i] = hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses{
			ContractAddress: ca,
		}
	}
	car := make([]hestia_autogen_bases.ProvisionedQuicknodeServicesReferers, len(pr.Referers))
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
	return nil
}

func (h *HestiaQuicknodeActivities) UpdateProvision(ctx context.Context, pr hestia_quicknode.ProvisionRequest) error {
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
	if pr.Referers == nil {
		pr.Referers = []string{}
	}
	if pr.ContractAddresses == nil {
		pr.ContractAddresses = []string{}
	}
	cas := make([]hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses, len(pr.ContractAddresses))
	for i, ca := range pr.ContractAddresses {
		cas[i] = hestia_autogen_bases.ProvisionedQuicknodeServicesContractAddresses{
			ContractAddress: ca,
		}
	}
	car := make([]hestia_autogen_bases.ProvisionedQuicknodeServicesReferers, len(pr.Referers))
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
		log.Warn().Err(err).Msg("Provision: UpdateProvision")
		return err
	}
	return nil
}

func (h *HestiaQuicknodeActivities) Deprovision(ctx context.Context, ou org_users.OrgUser, dp hestia_quicknode.DeprovisionRequest) error {
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

func (h *HestiaQuicknodeActivities) Deactivate(ctx context.Context, ou org_users.OrgUser, da hestia_quicknode.DeactivateRequest) error {
	err := hestia_quicknode_models.DeactivateProvisionedQuickNodeServiceEndpoint(ctx, da.QuickNodeID, da.EndpointID)
	if err != nil {
		log.Warn().Interface("ou", ou).Err(err).Msg("Provision: Deactivate")
		return err
	}
	return nil
}

func (h *HestiaQuicknodeActivities) CheckPlanOverages(ctx context.Context, pr hestia_quicknode.ProvisionRequest) ([]string, error) {
	if pr.QuickNodeID == "" {
		return nil, nil
	}
	tc, err := iris_models.OrgGroupTablesToRemove(context.Background(), pr.QuickNodeID, pr.Plan)
	if err != nil {
		log.Warn().Err(err).Msg("Provision: CheckPlanOverages")
		return nil, err
	}
	return tc, nil
}

func (h *HestiaQuicknodeActivities) DeactivateApiKey(ctx context.Context, ou org_users.OrgUser, pr hestia_quicknode.DeprovisionRequest) error {
	err := read_keys.DeactivateQuickNodeApiKey(context.Background(), pr.QuickNodeID)
	if err != nil {
		log.Warn().Msg("Provision: DeactivateApiKey")
		return err
	}
	return nil
}
