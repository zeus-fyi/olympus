package nodes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func ExternalApiNodesRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e.POST("/nodes", NodeActionsRequestHandler)
	return e
}

type ActionRequest struct {
	kns.TopologyKubeCtxNs
	Labels map[string]string `json:"labels,omitempty"`
	Action string            `json:"action"`
}

func NodeActionsRequestHandler(c echo.Context) error {
	request := new(ActionRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	if request.Action == "list" {
		return ListNodesRequest(c, request)
	}
	return c.JSON(http.StatusBadRequest, nil)
}
