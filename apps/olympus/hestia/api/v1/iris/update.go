package hestia_iris_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
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
	ou := c.Get("orgUser").(org_users.OrgUser)
	tc, err := iris_models.OrgEndpointsAndGroupTablesCount(context.Background(), ou.OrgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	token := c.Get("bearer").(string)
	if len(token) <= 0 {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	key, err := auth.VerifyBearerTokenServiceWithQuickNodePlan(context.Background(), token, create_org_users.IrisQuickNodeService)
	if err != nil {
		log.Err(err).Msg("InitV1Routes")
		return c.JSON(http.StatusUnauthorized, nil)
	}
	if len(key.PublicKeyName) <= 0 {
		return c.JSON(http.StatusUnauthorized, nil)
	}
	err = tc.CheckPlanLimits(key.PublicKeyName)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
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
	return c.JSON(http.StatusOK, nil)
}
