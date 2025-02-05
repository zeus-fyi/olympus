package v1

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
)

func Routes(e *echo.Echo) *echo.Echo {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.GET("/health", Health)

	v1Group := e.Group("/v1")
	v1Group.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(2)))
	v1Group.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "altDB", "auth")
			key, err := auth.VerifyBearerToken(ctx, token)
			return key.PublicKeyVerified, err
		},
	}))
	v1Group.POST("/validators", HandleValidatorsRequest)
	v1Group.POST("/validator_balances", HandleValidatorBalancesRequest)
	v1Group.POST("/validator_balances_sums", HandleValidatorBalancesSumRequest)

	debugGroup := e.Group("/debug")
	debugGroup.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == "hQyPerNFu7C9wMYpzTtZubP9BnUTzpCV5", nil
		},
	}))

	return e
}
