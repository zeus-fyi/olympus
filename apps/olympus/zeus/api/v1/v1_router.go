package zeus_v1_router

import (
	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	zeus_v1_ai "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/ai"
	zeus_v1_compute_api "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/compute"
	zeus_v1_clusters_api "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/ext_clusters"
	topology_routes "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func V1Routes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	topology_routes.Routes(e, k8Cfg)
	zeus_v1_compute_api.ComputeV1Routes(e)
	zeus_v1_ai.AiV1Routes(e)
	return e
}

func ExtSecureIntegrationRoutes(e *echo.Group) *echo.Group {
	zeus_v1_clusters_api.ExternalDeployRoutes(e)
	return e
}

func V1RoutesUI(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	topology_routes.RoutesUI(e, k8Cfg)
	return e
}
func ActionsV1Routes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	topology_routes.ActionsRoutes(e, k8Cfg)
	return e
}

func V1InternalRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	topology_routes.InternalRoutes(e, k8Cfg)
	return e
}
