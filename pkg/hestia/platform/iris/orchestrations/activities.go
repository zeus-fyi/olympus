package platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris/models/bases/autogen"
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

func (h *HestiaPlatformActivities) DeleteOrgGroupRoutingTable(ctx context.Context, pr IrisPlatformServiceRequest) error {
	err := iris_models.DeleteOrgGroupAndRoutes(context.Background(), pr.Ou.OrgID, pr.OrgGroupName)
	if err != nil {
		log.Err(err).Msg("DeleteOrgGroupRoutingTable")
		return err
	}
	return nil
}

func (h *HestiaPlatformActivities) DeleteOrgRoutes(ctx context.Context, pr IrisPlatformServiceRequest) error {
	err := iris_models.DeleteOrgRoutes(context.Background(), pr.Ou.OrgID, pr.Routes)
	if err != nil {
		log.Err(err).Msg("DeleteOrgGroupRoutingTable")
		return err
	}
	return nil
}

func (h *HestiaPlatformActivities) IrisPlatformSetupCacheUpdateRequest(ctx context.Context, pr IrisPlatformServiceRequest) error {
	// call iris via api & cache refresh

	return nil
}
