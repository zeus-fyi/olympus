package platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris/models/bases/autogen"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	"github.com/zeus-fyi/olympus/pkg/iris/resty_base"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
)

type HestiaPlatformActivities struct {
	kronos_helix.KronosActivities
}

func NewHestiaPlatformActivities() HestiaPlatformActivities {
	return HestiaPlatformActivities{
		KronosActivities: kronos_helix.NewKronosActivities(),
	}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (h *HestiaPlatformActivities) GetActivities() ActivitiesSlice {
	actSlice := []interface{}{
		h.IrisPlatformSetupCacheUpdateRequest, h.UpdateDatabaseOrgRoutingTables, h.CreateOrgGroupRoutingTable,
		h.DeleteOrgRoutes, h.IrisPlatformDeleteGroupTableCacheRequest,
		h.IrisPlatformDeleteOrgGroupTablesCacheRequest, h.DeleteOrgRoutesFromGroup,
		h.IrisPlatformRefreshOrgGroupTableCacheRequest, h.DeleteOrgRoutingGroup,
	}
	actSlice = append(actSlice, h.KronosActivities.GetActivities()...)
	return actSlice
}

func (h *HestiaPlatformActivities) DeleteOrgRoutes(ctx context.Context, pr IrisPlatformServiceRequest) error {
	if len(pr.Routes) == 0 {
		return nil
	}
	err := iris_models.DeleteOrgRoutes(context.Background(), pr.Ou.OrgID, pr.Routes)
	if err != nil {
		log.Err(err).Interface("pr", pr).Msg("HestiaPlatformActivities: DeleteOrgRoutes")
		return err
	}
	return nil
}

func (h *HestiaPlatformActivities) DeleteOrgRoutesFromGroup(ctx context.Context, pr IrisPlatformServiceRequest) error {
	if len(pr.Routes) == 0 {
		return nil
	}
	if pr.OrgGroupName == "" {
		return nil
	}

	err := iris_models.DeleteOrgRoutesFromGroup(context.Background(), pr.Ou.OrgID, pr.OrgGroupName, pr.Routes)
	if err != nil {
		log.Err(err).Interface("pr", pr).Msg("HestiaPlatformActivities: DeleteOrgRoutesFromGroup")
		return err
	}
	return nil
}

func (h *HestiaPlatformActivities) DeleteOrgRoutingGroup(ctx context.Context, pr IrisPlatformServiceRequest) error {
	if pr.OrgGroupName == "" {
		return nil
	}

	err := iris_models.DeleteOrgRoutingGroup(context.Background(), pr.Ou.OrgID, pr.OrgGroupName)
	if err != nil {
		log.Err(err).Interface("pr", pr).Msg("HestiaPlatformActivities: DeleteOrgRoutingGroup")
		return err
	}
	return nil
}

func (h *HestiaPlatformActivities) UpdateDatabaseOrgRoutingTables(ctx context.Context, pr IrisPlatformServiceRequest) error {
	routes := make([]iris_autogen_bases.OrgRoutes, len(pr.Routes))
	for i, route := range pr.Routes {
		routes[i] = iris_autogen_bases.OrgRoutes{
			RoutePath: route,
		}
	}
	err := iris_models.InsertOrgRoutes(context.Background(), pr.Ou.OrgID, routes)
	if err != nil {
		log.Err(err).Interface("pr", pr).Msg("HestiaPlatformActivities: UpdateDatabaseOrgRoutingTables")
		return err
	}
	return nil
}

func (h *HestiaPlatformActivities) CreateOrgGroupRoutingTable(ctx context.Context, pr IrisPlatformServiceRequest) error {
	if len(pr.Routes) == 0 {
		return nil
	}
	routes := make([]iris_autogen_bases.OrgRoutes, len(pr.Routes))
	for i, route := range pr.Routes {
		routes[i] = iris_autogen_bases.OrgRoutes{
			RoutePath: route,
		}
	}
	if pr.OrgGroupName == "" {
		return nil
	}
	ogr := iris_autogen_bases.OrgRouteGroups{
		OrgID:          pr.Ou.OrgID,
		RouteGroupName: pr.OrgGroupName,
	}
	err := iris_models.InsertOrgRouteGroup(context.Background(), ogr, routes)
	if err != nil {
		log.Err(err).Interface("ogr", ogr).Msg("InsertOrgRouteGroup")
		return err
	}
	return nil
}

var IrisApiUrl = "https://iris.zeus.fyi"

func (h *HestiaPlatformActivities) IrisPlatformSetupCacheUpdateRequest(ctx context.Context, pr IrisPlatformServiceRequest) error {
	rc := resty_base.GetBaseRestyClient(IrisApiUrl, artemis_orchestration_auth.Bearer)
	refreshEndpoint := fmt.Sprintf("/v1/internal/router/refresh/%d", pr.Ou.OrgID)
	resp, err := rc.R().Get(refreshEndpoint)
	if err != nil {
		log.Err(err).Msg("IrisPlatformSetupCacheUpdateRequest")
		return err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Interface("orgUser", pr.Ou).Msg("IrisPlatformSetupCacheUpdateRequest")
		return err
	}
	return nil
}

func (h *HestiaPlatformActivities) IrisPlatformDeleteGroupTableCacheRequest(ctx context.Context, pr IrisPlatformServiceRequest) error {
	rc := resty_base.GetBaseRestyClient(IrisApiUrl, artemis_orchestration_auth.Bearer)
	refreshEndpoint := fmt.Sprintf("/v1/internal/router/%d/%s", pr.Ou.OrgID, pr.OrgGroupName)
	resp, err := rc.R().Delete(refreshEndpoint)
	if err != nil {
		log.Err(err).Msg("IrisPlatformDeleteGroupTableCacheRequest")
		return err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Interface("orgUser", pr.Ou).Msg("HestiaPlatformActivities: IrisPlatformDeleteGroupTableCacheRequest")
		return err
	}
	return nil
}

func (h *HestiaPlatformActivities) IrisPlatformDeleteOrgGroupTablesCacheRequest(ctx context.Context, pr IrisPlatformServiceRequest) error {
	rc := resty_base.GetBaseRestyClient(IrisApiUrl, artemis_orchestration_auth.Bearer)
	refreshEndpoint := fmt.Sprintf("/v1/internal/router/%d", pr.Ou.OrgID)
	resp, err := rc.R().Delete(refreshEndpoint)
	if err != nil {
		log.Err(err).Msg("HestiaPlatformActivities: IrisPlatformDeleteOrgGroupTablesCacheRequest")
		return err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Interface("orgUser", pr.Ou).Msg("HestiaPlatformActivities: IrisPlatformDeleteOrgGroupTablesCacheRequest")
		return err
	}
	return nil
}

func (h *HestiaPlatformActivities) IrisPlatformRefreshOrgGroupTableCacheRequest(ctx context.Context, pr IrisPlatformServiceRequest) error {
	rc := resty_base.GetBaseRestyClient(IrisApiUrl, artemis_orchestration_auth.Bearer)
	refreshEndpoint := fmt.Sprintf("/v1/internal/router/refresh/%d/%s", pr.Ou.OrgID, pr.OrgGroupName)
	resp, err := rc.R().Get(refreshEndpoint)
	if err != nil {
		log.Err(err).Msg("HestiaPlatformActivities: IrisPlatformRefreshOrgGroupTableCacheRequest")
		return err
	}
	if resp.StatusCode() >= 400 {
		log.Err(err).Interface("orgUser", pr.Ou).Msg("HestiaPlatformActivities: IrisPlatformRefreshOrgGroupTableCacheRequest")
		return err
	}
	return nil
}
