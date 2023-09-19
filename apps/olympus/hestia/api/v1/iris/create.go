package hestia_iris_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	iris_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris/models/bases/autogen"
	platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/iris/orchestrations"
)

func CreateOrgRoutesRequestHandler(c echo.Context) error {
	request := new(OrgGroupRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request.GroupName != "" {
		return request.CreateGroupRoute(c)
	}
	return request.Create(c)
}

func (r *OrgGroupRoutesRequest) Create(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	tc, err := iris_models.OrgEndpointsAndGroupTablesCount(context.Background(), ou.OrgID, ou.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	err = tc.CheckEndpointLimits()
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}
	or := make([]iris_autogen_bases.OrgRoutes, len(r.Routes))
	for i, route := range r.Routes {
		or[i] = iris_autogen_bases.OrgRoutes{
			RoutePath: route,
		}
	}
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:     ou,
		Routes: r.Routes,
	}
	ctx := context.Background()
	err = platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisPlatformSetupRequestWorkflow(ctx, ipr)
	if err != nil {
		log.Err(err).Msg("CreateOrgGroupRoutesRequest")
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

const QuickNodeMarketPlace = "quickNodeMarketPlace"

func CreateOrgGroupRoutesRequestHandler(c echo.Context) error {
	request := new(OrgGroupRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateGroupRoute(c)
}

type OrgGroupRoutesRequest struct {
	GroupName string   `json:"groupName,omitempty"`
	Routes    []string `json:"routes,omitempty"`
}

func (r *OrgGroupRoutesRequest) CreateGroupRoute(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	sp, ok := c.Get("servicePlans").(map[string]string)
	if !ok {
		log.Warn().Interface("servicePlans", sp).Msg("CreateGroupRoute: marketplace not found")
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}
	plan, ok := sp[QuickNodeMarketPlace]
	if !ok {
		log.Warn().Str("marketplace", QuickNodeMarketPlace).Msg("CreateGroupRoute: marketplace not found")
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}
	tc, err := iris_models.OrgEndpointsAndGroupTablesCount(context.Background(), ou.OrgID, ou.UserID)
	if err != nil {
		log.Err(err).Msg("CreateGroupRoute: OrgEndpointsAndGroupTablesCount")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	err = tc.CheckPlanLimits(plan)
	if err != nil {
		log.Err(err).Interface("plan", plan).Msg("CreateGroupRoute: CheckPlanLimits")
		return c.JSON(http.StatusPreconditionFailed, err)
	}
	or := make([]iris_autogen_bases.OrgRoutes, len(r.Routes))
	for i, route := range r.Routes {
		or[i] = iris_autogen_bases.OrgRoutes{
			RoutePath: route,
		}
	}
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:           ou,
		OrgGroupName: r.GroupName,
		Routes:       r.Routes,
	}
	ctx := context.Background()
	err = platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisPlatformSetupRequestWorkflow(ctx, ipr)
	if err != nil {
		log.Err(err).Interface("ou", ou).Msg("ExecuteIrisPlatformSetupRequestWorkflow")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}
