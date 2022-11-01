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
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/core"
)

func InitRouter(e *echo.Echo, k8Cfg autok8s_core.K8Util) *echo.Echo {
	log.Debug().Msgf("InitRouter")
	k8Cfg.ConnectToK8sFromConfig(k8Cfg.CfgPath)
	e = Routes(e, k8Cfg)
	return e
}

func Routes(e *echo.Echo, k8Cfg autok8s_core.K8Util) *echo.Echo {
	core.K8Util = k8Cfg

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.GET("/health", Health)

	v1RoutesGroup := v1.V1Routes(e, k8Cfg)
	v1RoutesGroup.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Bearer",
		Validator: func(token string, c echo.Context) (bool, error) {
			ctx := context.Background()
			key, err := auth.VerifyBearerToken(ctx, token)
			return key.PublicKeyVerified, err
		},
	}))
	return e
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
