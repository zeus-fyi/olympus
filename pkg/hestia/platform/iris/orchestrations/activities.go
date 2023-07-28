package platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris/models/bases/autogen"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	"github.com/zeus-fyi/olympus/pkg/iris/resty_base"
)

type HestiaPlatformActivities struct {
}

func NewHestiaPlatformActivities() HestiaPlatformActivities {
	return HestiaPlatformActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (h *HestiaPlatformActivities) GetActivities() ActivitiesSlice {
	return []interface{}{
		h.IrisPlatformSetupCacheUpdateRequest, h.UpdateDatabaseOrgRoutingTables, h.CreateOrgGroupRoutingTable,
		h.DeleteOrgGroupRoutingTable, h.DeleteOrgRoutes,
	}
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
		log.Err(err).Msg("UpdateDatabaseOrgRoutingTables")
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
		log.Err(err).Msg("UpdateDatabaseOrgRoutingTables")
		return err
	}
	return nil
}

func (h *HestiaPlatformActivities) DeleteOrgGroupRoutingTable(ctx context.Context, pr IrisPlatformServiceRequest) error {
	if pr.OrgGroupName == "" {
		return nil
	}
	err := iris_models.DeleteOrgGroupAndRoutes(context.Background(), pr.Ou.OrgID, pr.OrgGroupName)
	if err != nil {
		log.Err(err).Msg("DeleteOrgGroupRoutingTable")
		return err
	}
	return nil
}

func (h *HestiaPlatformActivities) DeleteOrgRoutes(ctx context.Context, pr IrisPlatformServiceRequest) error {
	if len(pr.Routes) == 0 {
		return nil
	}
	err := iris_models.DeleteOrgRoutes(context.Background(), pr.Ou.OrgID, pr.Routes)
	if err != nil {
		log.Err(err).Msg("DeleteOrgGroupRoutingTable")
		return err
	}
	return nil
}

const (
	IrisApiUrl = "https://iris.zeus.fyi"
)

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
