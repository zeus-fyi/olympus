package topology_routes

import (
	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	topology_auths "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/auth"
	deploy_routes "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy"
	external_actions_routes "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/external_actions"
	infra_routes "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func Routes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	topology_auths.K8Util = k8Cfg
	e = infra_routes.Routes(e, k8Cfg)
	e = deploy_routes.ExternalDeployRoutes(e, k8Cfg)
	return e
}

func RoutesUI(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	topology_auths.K8Util = k8Cfg
	e = infra_routes.UIRoutes(e, k8Cfg)
	return e
}

func ActionsRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	topology_auths.K8Util = k8Cfg
	e = external_actions_routes.ExternalActionsRoutes(e, k8Cfg)
	return e
}

func InternalRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e = deploy_routes.InternalRoutes(e, k8Cfg)
	return e
}
