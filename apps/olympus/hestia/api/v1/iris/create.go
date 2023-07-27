package hestia_iris_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris/models/bases/autogen"
	platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/iris/orchestrations"
)

func CreateOrgRoutesRequestHandler(c echo.Context) error {
	request := new(CreateOrgRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Create(c)
}

type CreateOrgRoutesRequest struct {
	Routes []string `json:"routes"`
}

func (r *CreateOrgRoutesRequest) Create(c echo.Context) error {
	or := make([]iris_autogen_bases.OrgRoutes, len(r.Routes))
	for i, route := range r.Routes {
		or[i] = iris_autogen_bases.OrgRoutes{
			RoutePath: route,
		}
	}
	ou := c.Get("orgUser").(org_users.OrgUser)
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:     ou,
		Routes: r.Routes,
	}
	ctx := context.Background()
	err := platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisPlatformSetupRequestWorkflow(ctx, ipr)
	if err != nil {
		log.Err(err).Msg("CreateOrgGroupRoutesRequest")
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

func CreateOrgGroupRoutesRequestHandler(c echo.Context) error {
	request := new(CreateOrgRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Create(c)
}

type CreateOrUpdateOrgGroupRoutesRequest struct {
	GroupName string   `json:"groupName"`
	Routes    []string `json:"routes"`
}

func (r *CreateOrUpdateOrgGroupRoutesRequest) Create(c echo.Context) error {
	or := make([]iris_autogen_bases.OrgRoutes, len(r.Routes))
	for i, route := range r.Routes {
		or[i] = iris_autogen_bases.OrgRoutes{
			RoutePath: route,
		}
	}
	ou := c.Get("orgUser").(org_users.OrgUser)
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:           ou,
		OrgGroupName: r.GroupName,
		Routes:       r.Routes,
	}
	ctx := context.Background()
	err := platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisPlatformSetupRequestWorkflow(ctx, ipr)
	if err != nil {
		log.Err(err).Msg("CreateOrgGroupRoutesRequest")
		return err
	}
	return c.JSON(http.StatusOK, nil)
}
