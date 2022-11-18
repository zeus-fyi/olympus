package athena_router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	v1_athena "github.com/zeus-fyi/olympus/athena/api/v1"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func Routes(e *echo.Echo, p filepaths.Path) *echo.Echo {
	// Routes
	e.GET("/health", Health)

	v1_athena.InitV1InternalRoutes(e, p)
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
