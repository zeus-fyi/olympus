package router

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	v1 "github.com/zeus-fyi/olympus/zeus/api/v1"
)

func InitRouter(e *echo.Echo, k8Cfg autok8s_core.K8Util) *echo.Echo {
	log.Debug().Msgf("InitRouter")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.GET("/health", Health)

	InitV1Routes(e, k8Cfg)
	return e
}

func InitV1Routes(e *echo.Echo, k8Cfg autok8s_core.K8Util) {
	eg := e.Group("/v1")
	eg = v1.V1Routes(eg, k8Cfg)
	eg.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			key, err := auth.VerifyBearerToken(ctx, token)
			return key.PublicKeyVerified, err
		},
	}))
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
