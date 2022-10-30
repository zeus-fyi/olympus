package coreK8s

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

var K8util autok8s_core.K8Util

type K8sRequest struct {
	Kns autok8s_core.KubeCtxNs
}

func InitRouter(e *echo.Echo, k8Cfg autok8s_core.K8Util) *echo.Echo {
	log.Debug().Msgf("InitRouter")
	k8Cfg.ConnectToK8sFromConfig(k8Cfg.CfgPath)
	e = Routes(e, k8Cfg)
	return e
}

func Routes(e *echo.Echo, k8Cfg autok8s_core.K8Util) *echo.Echo {
	K8util = k8Cfg
	// TODO add authentication
	e.POST("/pods", HandlePodActionRequest)
	e.POST("/topology", HandleTopologyActionRequest)

	return e
}
