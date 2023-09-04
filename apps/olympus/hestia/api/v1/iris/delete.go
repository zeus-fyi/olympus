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
	if len(r.Routes) == 0 && r.GroupName == "" {
		return c.JSON(http.StatusBadRequest, "no routes provided for deletion")
	}
	if len(r.Routes) == 0 && len(r.GroupName) > 0 {
		return r.DeleteOrgRoutingGroup(c)
	}
	ou := c.Get("orgUser").(org_users.OrgUser)
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:     ou,
		Routes: r.Routes,
	}
	err := platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisDeleteOrgRoutesWorkflow(context.Background(), ipr)
	if err != nil {
		log.Err(err).Msg("DeleteOrgRoutes")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, QuickNodeResponse{
		Status: "success",
	})
}

func InternalDeleteOrgRoutesRequestHandler(c echo.Context) error {
	request := new(OrgGroupRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.InternalDeleteOrgRoutes(c)
}

func (r *OrgGroupRoutesRequest) InternalDeleteOrgRoutes(c echo.Context) error {
	if len(r.Routes) == 0 {
		return c.JSON(http.StatusBadRequest, "no routes provided for deletion")
	}
	orgID, ok := c.Get("orgID").(int)
	if !ok {
		return c.JSON(http.StatusBadRequest, "orgID is required")
	}
	ou := org_users.NewOrgUserWithID(orgID, 0)
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:     ou,
		Routes: r.Routes,
	}
	err := platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisDeleteOrgRoutesWorkflow(context.Background(), ipr)
	if err != nil {
		log.Err(err).Msg("DeleteOrgRoutes")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, QuickNodeResponse{
		Status: "success",
	})
}

func DeleteReplaceOrgGroupRoutesRequestHandler(c echo.Context) error {
	request := new(OrgGroupRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.DeleteOrgRoutingGroup(c)
}

func (r *OrgGroupRoutesRequest) DeleteOrgRoutingGroup(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusBadRequest, nil)
	}
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
	return c.JSON(http.StatusOK, QuickNodeResponse{
		Status: "success",
	})
}

func DeletePartialOrgGroupRoutesRequestHandler(c echo.Context) error {
	request := new(OrgGroupRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.DeletePartialOrgGroupRoutes(c)
}

func (r *OrgGroupRoutesRequest) DeletePartialOrgGroupRoutes(c echo.Context) error {
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if len(r.GroupName) == 0 {
		return c.JSON(http.StatusBadRequest, "GroupName is required")
	}
	if (len(r.Routes)) == 0 {
		return c.JSON(http.StatusBadRequest, "Routes are required")
	}
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:           ou,
		Routes:       r.Routes,
		OrgGroupName: r.GroupName,
	}
	err := platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisRemoveRoutesFromOrgGroupRoutingTableWorkflow(context.Background(), ipr)
	if err != nil {
		log.Err(err).Msg("DeletePartialOrgGroupRoutes: ExecuteIrisRemoveRoutesFromOrgGroupRoutingTableWorkflow")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, QuickNodeResponse{
		Status: "success",
	})
}
