package infra

import (
	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func Routes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	e.POST("/infra", HandleTopologyInfraActionRequest)
	return e
}
