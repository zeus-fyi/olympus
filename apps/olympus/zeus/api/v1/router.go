package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology"
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
	// TODO add authentication

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// Routes
	e.GET("/health", Health)

	topologyV1Routes := topology.Routes(e, k8Cfg)
	v1TopologyAPIGroup := topologyV1Routes.Group("/v1")
	v1TopologyAPIGroup.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
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
