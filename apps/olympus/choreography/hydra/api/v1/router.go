package v1_hydra_choreography

import (
	"net/http"

	"github.com/labstack/echo/v4"
	hydra_choreography_metrics "github.com/zeus-fyi/olympus/choreography/hydra/metrics"
	zeus_client "github.com/zeus-fyi/zeus/zeus/z_client"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var (
	ZeusClient zeus_client.ZeusClient
	CloudCtxNs zeus_common_types.CloudCtxNs
)

func Routes(e *echo.Echo) *echo.Echo {
	e.GET("/health", Health)
	e.GET("/delete/pods", RestartPods)
	e.GET("/metrics", hydra_choreography_metrics.MetricsRequestHandler)
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
