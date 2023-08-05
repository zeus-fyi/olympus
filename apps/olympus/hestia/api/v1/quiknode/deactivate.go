package hestia_quiknode_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode"
	quicknode_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode/orchestrations"
)

func DeactivateRequestHandler(c echo.Context) error {
	request := new(DeactivateRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Deactivate(c)
}

type DeactivateRequest struct {
	hestia_quicknode.DeactivateRequest
}

func (r *DeactivateRequest) Deactivate(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	da := r.DeactivateRequest
	err := quicknode_orchestrations.HestiaQnWorker.ExecuteQnDeactivateWorkflow(context.Background(), ou, da)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			QuickNodeResponse{
				Status: "error",
			})
	}
	return c.JSON(http.StatusOK, QuickNodeResponse{
		Status: "success",
	})
}
