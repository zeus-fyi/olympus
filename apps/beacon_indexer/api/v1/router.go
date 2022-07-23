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

	e.GET("/log/:level", SetLogLevel)
	e.GET("/validator/new/:batchSize", SetNewValidatorBatchSize)
	e.GET("/validator/balances/:batchSize", SetNewValidatorBalanceBatchSize)
	return e
}
