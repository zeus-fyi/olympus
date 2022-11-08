package create_infra

import (
	"github.com/labstack/echo/v4"
)

func CreateTopologyInfraActionRequestHandler(c echo.Context) error {
	request := new(TopologyCreateRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	return request.CreateTopology(c)
}
