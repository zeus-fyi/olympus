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

func CreateTopologyInfraActionFromUIRequestHandler(c echo.Context) error {
	request := new(TopologyCreateRequestFromUI)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateTopologyFromUI(c)
}

func CreateTopologyClassActionRequestHandler(c echo.Context) error {
	request := new(TopologyCreateOrAddBasesToClassesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.CreateTopologyClusterClass(c)
}

func UpdateTopologyClassActionRequestHandler(c echo.Context) error {
	request := new(TopologyCreateOrAddComponentBasesToClassesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.AddComponentBasesToTopologyClusterClass(c)
}

func CreateTopologySkeletonBasesActionRequestHandler(c echo.Context) error {
	request := new(TopologyCreateOrAddBasesToClassesRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.AddSkeletonBaseClassToBase(c)
}
