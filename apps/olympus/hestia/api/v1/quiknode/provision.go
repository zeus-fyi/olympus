package hestia_quiknode_v1_routes

import (
	"github.com/labstack/echo/v4"
	hestia_quiknode "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/quiknode"
)

func ProvisionRequestHandler(c echo.Context) error {
	request := new(ProvisionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Provision(c)
}

type ProvisionRequest struct {
	hestia_quiknode.ProvisionRequest
}

func (r *ProvisionRequest) Provision(c echo.Context) error {
	return nil
}
