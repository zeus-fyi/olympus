package hestia_quiknode_v1_routes

import (
	"github.com/labstack/echo/v4"
	hestia_quiknode "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/quiknode"
)

func DeprovisionRequestHandler(c echo.Context) error {
	request := new(DeprovisionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Deprovision(c)
}

type DeprovisionRequest struct {
	hestia_quiknode.DeprovisionRequest
}

func (r *DeprovisionRequest) Deprovision(c echo.Context) error {
	return nil
}
