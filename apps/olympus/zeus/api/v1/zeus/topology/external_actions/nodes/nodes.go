package nodes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
)

type ActionRequest struct {
	kns.TopologyKubeCtxNs
	Action string `json:"action"`
}

func NodeActionsRequestHandler(c echo.Context) error {
	request := c.Get("NodeActionsRequestHandler").(*ActionRequest)

	if request.Action == "list" {
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusBadRequest, nil)
}
