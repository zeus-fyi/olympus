package nodes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

func ExternalApiNodesRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e.POST("/nodes", NodeActionsRequestHandler)
	return e
}

type ActionRequest struct {
	Labels                       map[string]string `json:"labels,omitempty"`
	Action                       string            `json:"action"`
	zeus_common_types.CloudCtxNs `json:"cloudCtx,omitempty"`
}

func NodeActionsRequestHandler(c echo.Context) error {
	log.Info().Msg("NodeActionsRequestHandler")
	request := new(ActionRequest)
	if err := c.Bind(request); err != nil {
		log.Err(err).Msg("NodeActionsRequestHandler, Bind error")
		return err
	}
	if request.Action == "list" {
		return ListNodesRequest(c, request)
	}
	log.Info().Msg("NodeActionsRequestHandler, no valid actions found")
	return c.JSON(http.StatusBadRequest, nil)
}
