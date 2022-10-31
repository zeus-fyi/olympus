package infra

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleTopologyInfraActionRequest(c echo.Context) error {
	request := new(TopologyInfraRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request.Action == "create" {
		return request.CreateTopology(c)
	}
	if request.Action == "read" {
		return request.ReadTopology(c)
	}
	if request.Action == "update" {
		return request.UpdateTopology(c)
	}
	if request.Action == "delete" {
		return request.DeleteTopology(c)
	}
	return c.JSON(http.StatusBadRequest, nil)
}
