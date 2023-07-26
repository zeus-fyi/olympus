package hestia_quiknode_v1_routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/quiknode"
	quicknode_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/quiknode/orchestrations"
)

func ProvisionRequestHandler(c echo.Context) error {
	request := new(ProvisionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Provision(c)
}

type ProvisionRequest struct {
	hestia_quicknode.ProvisionRequest
}

func (r *ProvisionRequest) Provision(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	pr := r.ProvisionRequest
	err := quicknode_orchestrations.HestiaQnWorker.ExecuteQnProvisionWorkflow(context.Background(), pr, ou)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "success")
}

func UpdateProvisionRequestHandler(c echo.Context) error {
	request := new(ProvisionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.UpdateProvision(c)
}

func (r *ProvisionRequest) UpdateProvision(c echo.Context) error {
	ou := c.Get("orgUser").(org_users.OrgUser)
	pr := r.ProvisionRequest
	err := quicknode_orchestrations.HestiaQnWorker.ExecuteQnUpdateProvisionWorkflow(context.Background(), pr, ou)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "success")
}
