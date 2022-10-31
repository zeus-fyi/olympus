package pods

import (
	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus/zeus_pkg"
)

func Routes(e *echo.Echo, k8Cfg autok8s_core.K8Util) *echo.Echo {
	zeus_pkg.K8Util = k8Cfg
	// TODO add authentication
	e.POST("/pods", HandlePodActionRequest)

	return e
}
