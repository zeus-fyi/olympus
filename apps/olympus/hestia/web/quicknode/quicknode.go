package hestia_quicknode_dashboard

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var JWTAuthSecret = ""

func InitQuickNodeDashboardRoutes(e *echo.Echo) {
	eg := e.Group("/v1/quicknode/dashboard")
	eg.Use(echojwt.JWT(middleware.JWTConfig{
		SigningKey: []byte(JWTAuthSecret),
	}))
}
