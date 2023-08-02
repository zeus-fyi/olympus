package hestia_quiknode_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode"
	quicknode_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode/orchestrations"
)

func DeprovisionRequestHandler(c echo.Context) error {
	request := new(DeprovisionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Deprovision(c)
}

type DeprovisionRequest struct {
	hestia_quicknode.DeprovisionRequest
}

func (r *DeprovisionRequest) Deprovision(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	dp := r.DeprovisionRequest
	err := quicknode_orchestrations.HestiaQnWorker.ExecuteQnDeprovisionWorkflow(context.Background(), ou, dp)
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
