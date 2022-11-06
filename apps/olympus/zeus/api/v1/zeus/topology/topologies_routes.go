package topology

import (
	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/actions"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func Routes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	e = actions.Routes(e, k8Cfg)
	e = infra.Routes(e, k8Cfg)
	e = deploy.Routes(e, k8Cfg)
	return e
}

func InternalRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e = deploy.InternalRoutes(e, k8Cfg)
	return e
}
