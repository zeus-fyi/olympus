package v1_hydra_choreography

import (
	"net/http"

	"github.com/labstack/echo/v4"
	zeus_client "github.com/zeus-fyi/zeus/pkg/zeus/client"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

var (
	ZeusClient zeus_client.ZeusClient
	CloudCtxNs zeus_common_types.CloudCtxNs
)

func Routes(e *echo.Echo) *echo.Echo {
	e.GET("/health", Health)
	e.GET("/delete/pods", RestartPods)
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
