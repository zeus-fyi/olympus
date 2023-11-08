package zeus_v1_compute_api

import (
	"github.com/labstack/echo/v4"
)

func ComputeV1Routes(e *echo.Group) *echo.Group {
	e.POST("/search/nodes", NodeSearchHandler)
	return e
}
