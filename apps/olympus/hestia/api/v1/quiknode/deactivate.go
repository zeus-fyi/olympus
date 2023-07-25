package hestia_quiknode_v1_routes

import (
	"github.com/labstack/echo/v4"
	hestia_quiknode "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/quiknode"
)

func DeactivateRequestHandler(c echo.Context) error {
	request := new(DeprovisionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Deprovision(c)
}

type DeactivateRequest struct {
	hestia_quiknode.DeactivateRequest
}

func (r *DeactivateRequest) Deactivate(c echo.Context) error {
	return nil
}
