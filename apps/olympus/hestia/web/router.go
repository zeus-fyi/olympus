package hestia_web_router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	hestia_login "github.com/zeus-fyi/olympus/hestia/web/login"
)

func WebRoutes(e *echo.Echo) *echo.Echo {
	// Routes

	e.POST("/login", hestia_login.CreateLoginHandler)
	return e
}

func Login(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
