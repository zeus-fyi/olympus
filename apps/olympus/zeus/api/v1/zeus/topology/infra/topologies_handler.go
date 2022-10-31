package infra

import (
	"net/http"

	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create"
	delete_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/delete"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/read"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/update"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/core"
)

type TopologyInfraRequest struct {
	Action string

	create_infra.TopologyActionCreateRequest
	read_infra.TopologyActionReadRequest
	update_infra.TopologyActionUpdateRequest
	delete_infra.TopologyActionDeleteRequest
}

func Routes(e *echo.Echo, k8Cfg autok8s_core.K8Util) *echo.Echo {
	core.K8Util = k8Cfg
	e.POST("/infra", HandleTopologyInfraActionRequest)
	return e
}

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
