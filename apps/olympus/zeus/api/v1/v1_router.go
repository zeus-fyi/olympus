package v1_router

import (
	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func V1Routes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	topology.Routes(e, k8Cfg)
	return e
}

func V1InternalRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	topology.InternalRoutes(e, k8Cfg)
	return e
}
