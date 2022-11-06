package deploy

import (
	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/internal/deploy/workload_deploy"
	internal_destroy_deploy "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/internal/deploy/workload_destroy"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func Routes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e.POST("/deploy", HandleDeploymentActionRequest)
	return e
}

func InternalRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	e = InternalDeployRoutes(e, k8Cfg)
	e = InternalDeployDestroyRoutes(e, k8Cfg)
	return e
}

func InternalDeployRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e.Group("/deploy")
	e.POST("/namespace", internal_deploy.DeployNamespaceHandler)
	e.POST("/deployment", internal_deploy.DeployDeploymentHandler)
	e.POST("/statefulset", internal_deploy.DeployStatefulSetHandler)
	e.POST("/configmap", internal_deploy.DeployConfigMapHandler)
	e.POST("/service", internal_deploy.DeployServiceHandler)
	e.POST("/ingress", internal_deploy.DeployIngressHandler)
	return e
}

func InternalDeployDestroyRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e.Group("/deploy/destroy")
	e.POST("/namespace", internal_destroy_deploy.DestroyDeployNamespaceHandler)
	e.POST("/deployment", internal_destroy_deploy.DestroyDeployDeploymentHandler)
	e.POST("/statefulset", internal_destroy_deploy.DestroyDeployStatefulSetHandler)
	e.POST("/configmap", internal_destroy_deploy.DestroyDeployConfigMapHandler)
	e.POST("/service", internal_destroy_deploy.DestroyDeployServiceHandler)
	e.POST("/ingress", internal_destroy_deploy.DestroyDeployIngressHandler)
	return e
}
