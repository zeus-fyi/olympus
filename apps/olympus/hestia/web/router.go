package hestia_web_router

import (
	"github.com/labstack/echo/v4"
	hestia_login "github.com/zeus-fyi/olympus/hestia/web/login"
)

func WebRoutes(e *echo.Echo) *echo.Echo {
	e.POST("/login", hestia_login.LoginHandler)
	return e
}
