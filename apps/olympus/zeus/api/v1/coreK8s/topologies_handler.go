package coreK8s

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleTopologyActionRequest(c echo.Context) error {
	request := new(TopologyActionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request.Action == "create" {
		return request.CreateTopology(c, request)
	}
	if request.Action == "deploy" {
		return request.DeployTopology(c, request)
	}
	if request.Action == "read" {
		return request.ReadTopology(c, request)
	}
	if request.Action == "update" {
		return request.UpdateTopology(c, request)
	}
	if request.Action == "update" {
		return request.DeleteTopology(c, request)
	}
	return c.JSON(http.StatusBadRequest, nil)
}
