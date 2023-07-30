package hestia_iris_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/iris/orchestrations"
)

func DeleteOrgRoutesRequestHandler(c echo.Context) error {
	request := new(OrgGroupRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.DeleteOrgRoutes(c)
}

func (r *OrgGroupRoutesRequest) DeleteOrgRoutes(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	if len(r.Routes) == 0 {
		return c.JSON(http.StatusBadRequest, "no routes provided for deletion")
	}
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:     ou,
		Routes: r.Routes,
	}
	err := platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisDeleteOrgRoutesWorkflow(context.Background(), ipr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return nil
}

func DeleteOrgGroupRoutesRequestHandler(c echo.Context) error {
	request := new(OrgGroupRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.DeleteOrgRoutingGroup(c)
}

func (r *OrgGroupRoutesRequest) DeleteOrgRoutingGroup(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	if len(r.GroupName) == 0 {
		return c.JSON(http.StatusBadRequest, "GroupName is required")
	}
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:           ou,
		OrgGroupName: r.GroupName,
	}
	err := platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisDeleteOrgGroupRoutingTableWorkflow(context.Background(), ipr)
	if err != nil {
		log.Err(err).Msg("DeleteOrgRoutingGroup")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}
