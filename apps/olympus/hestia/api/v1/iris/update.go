package hestia_iris_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
	hestia_quiknode_v1_routes "github.com/zeus-fyi/olympus/hestia/api/v1/quiknode"
	platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/iris/orchestrations"
)

func UpdateOrgGroupRoutesRequestHandler(c echo.Context) error {
	request := new(OrgGroupRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.UpdateOrgGroup(c)
}

func (r *OrgGroupRoutesRequest) UpdateOrgGroup(c echo.Context) error {
	if len(r.GroupName) == 0 {
		return c.JSON(http.StatusBadRequest, "GroupName is required")
	}
	if len(r.Routes) <= 0 {
		return r.DeleteOrgRoutingGroup(c)
	}
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	tc, err := iris_models.OrgEndpointsAndGroupTablesCount(context.Background(), ou.OrgID, ou.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
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
	err = tc.CheckPlanLimits(plan)
	if err != nil {
		return c.JSON(http.StatusPreconditionFailed, err)
	}
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:           ou,
		Routes:       r.Routes,
		OrgGroupName: r.GroupName,
	}
	err = platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisPlatformSetupRequestWorkflow(context.Background(), ipr)
	if err != nil {
		log.Err(err).Msg("UpdateOrgGroupRoutesRequest")
		return err
	}

	return c.JSON(http.StatusOK, hestia_quiknode_v1_routes.QuickNodeResponse{
		Status: "success",
	})
}

type QuickNodeResponse struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}
