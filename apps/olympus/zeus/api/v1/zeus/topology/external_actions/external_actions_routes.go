package external_actions_routes

import (
	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/external_actions/pods"
	external_api_workloads "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/external_actions/workloads"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func ExternalActionsRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	e = external_api_workloads.ExternalApiWorkloadQueryActionRoutes(e, k8Cfg)
	e = pods.ExternalApiPodsRoutes(e, k8Cfg)
	return e
}
