package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/pods"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/zeus_pkg"
)

func InitRouter(e *echo.Echo, k8Cfg autok8s_core.K8Util) *echo.Echo {
	log.Debug().Msgf("InitRouter")
	k8Cfg.ConnectToK8sFromConfig(k8Cfg.CfgPath)
	e = Routes(e, k8Cfg)
	return e
}

func Routes(e *echo.Echo, k8Cfg autok8s_core.K8Util) *echo.Echo {
	zeus_pkg.K8Util = k8Cfg
	// TODO add authentication
	e.POST("/pods", pods.HandlePodActionRequest)
	e.POST("/topology", topology.HandleTopologyActionRequest)

	return e
}
