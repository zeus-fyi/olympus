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

func CreateTopologyClassActionRequestHandler(c echo.Context) error {
	request := new(TopologyCreateClusterRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateTopologyClusterClass(c)
}

func UpdateTopologyClassActionRequestHandler(c echo.Context) error {
	request := new(TopologyAddBasesToClusterRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.AddBasesToTopologyClusterClass(c)
}
