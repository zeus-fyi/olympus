package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.GET("/health", Health)

	// TODO add users/orgs
	v1Group := e.Group("/v1")
	v1Group.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(2)))
	v1Group.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == "hQyPerNFu7C9wMYpzTtZubP9BnUTzpCV5", nil
		},
	}))

	debugGroup := e.Group("/debug")
	debugGroup.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == "hQyPerNFu7C9wMYpzTtZubP9BnUTzpCV5", nil
		},
	}))

	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
