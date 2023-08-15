package iris_api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	iris_metrics "github.com/zeus-fyi/olympus/iris/api/metrics"
	v1_iris "github.com/zeus-fyi/olympus/iris/api/v1"
	v1Beta_iris "github.com/zeus-fyi/olympus/iris/api/v1beta"
	v1internal_iris "github.com/zeus-fyi/olympus/iris/api/v1internal"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Routes
	e.GET("/healthz", Health)
	e.GET("/healthcheck", Health)
	e.GET("/health", Health)

	v1_iris.InitV1Routes(e)
	v1internal_iris.InitV1InternalRoutes(e)
	v1Beta_iris.InitV1BetaInternalRoutes(e)
	return e
}

type Response struct {
	Message string `json:"message"`
}

func Health(c echo.Context) error {
	resp := Response{Message: "ok"}
	return c.JSON(http.StatusOK, resp)
}

func MetricRoutes(e *echo.Echo) *echo.Echo {
	e.GET("/metrics", iris_metrics.MetricsRequestHandler)
	return e
}
