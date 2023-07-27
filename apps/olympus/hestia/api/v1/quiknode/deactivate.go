package hestia_quiknode_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/quiknode"
	quicknode_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/quiknode/orchestrations"
)

func DeactivateRequestHandler(c echo.Context) error {
	request := new(DeprovisionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Deprovision(c)
}

type DeactivateRequest struct {
	hestia_quicknode.DeactivateRequest
}

func (r *DeactivateRequest) Deactivate(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	da := r.DeactivateRequest
	err := quicknode_orchestrations.HestiaQnWorker.ExecuteQnDeactivateWorkflow(context.Background(), da, ou)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "success")
}
