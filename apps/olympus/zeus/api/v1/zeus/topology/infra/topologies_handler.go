package infra

import (
	"net/http"

	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/read"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/update"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/core"
)

type ActionRequest struct {
	Action string

	create.TopologyActionCreateRequest
	read.TopologyActionReadRequest
	update.TopologyActionUpdateRequest
}

func Routes(e *echo.Echo, k8Cfg autok8s_core.K8Util) *echo.Echo {
	core.K8Util = k8Cfg
	e.POST("/topology/infra", HandleTopologyInfraActionRequest)

	return e
}

func HandleTopologyInfraActionRequest(c echo.Context) error {
	request := new(ActionRequest)
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
	return c.JSON(http.StatusBadRequest, nil)
}
