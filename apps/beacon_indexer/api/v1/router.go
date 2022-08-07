package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.GET("/health", Health)
	//e.POST("/admin", HandleAdminConfigRequest)
	//e.GET("/admin", AdminGetRequestHandler)

	//e.GET("/debug/redis", DebugReadRedisRequestHandler)
	//e.POST("/debug/redis", DebugRedisRequestHandler)

	e.GET("/debug/db/counts", DebugRequestHandler)
	e.GET("/debug/db/sizes", TableSizesHandler)
	e.GET("/debug/db/stats", DebugPgStatsHandler)
	//e.GET("/debug/db/ping", PingDBHandler)
	//e.GET("/debug/db/config", DebugGetPgConfigHandler)

	//e.POST("/debug/db/config", DebugUpdatePgConfigHandler)
	return e
}
