package hestia_iris_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/iris/orchestrations"
)

func UpdateOrgGroupRoutesRequestHandler(c echo.Context) error {
	request := new(UpdateOrgGroupRoutesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Update(c)
}

type UpdateOrgGroupRoutesRequest struct {
	GroupName string   `json:"groupName,omitempty"`
	Routes    []string `json:"routes"`
}

func (r *UpdateOrgGroupRoutesRequest) Update(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	ipr := platform_service_orchestrations.IrisPlatformServiceRequest{
		Ou:           ou,
		Routes:       r.Routes,
		OrgGroupName: r.GroupName,
	}
	err := platform_service_orchestrations.HestiaPlatformServiceWorker.ExecuteIrisPlatformSetupRequestWorkflow(context.Background(), ipr)
	if err != nil {
		log.Err(err).Msg("UpdateOrgGroupRoutesRequest")
		return err
	}
	return c.JSON(http.StatusOK, nil)
}
