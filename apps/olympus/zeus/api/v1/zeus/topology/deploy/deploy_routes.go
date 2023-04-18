package deploy_routes

import (
	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	clean_deploy_request "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/clean"
	create_or_update_deploy "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/create_or_update"
	destroy_deploy_request "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/destroy"
	deployment_status "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/read"
	replace_topology "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/replace"
	internal_secrets_deploy "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/deploy/secrets_deploy"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/deploy/workload_deploy"
	internal_destroy_deploy "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/deploy/workload_destroy"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/deploy/workload_state"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func ExternalDeployRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	e.POST("/deploy", create_or_update_deploy.TopologyDeploymentHandler)
	e.POST("/deploy/cluster", create_or_update_deploy.ClusterTopologyDeploymentHandler)
	e.POST("/deploy", create_or_update_deploy.TopologyDeploymentHandler)
	e.POST("/deploy/ui/app", create_or_update_deploy.SetupClusterTopologyDeploymentHandler)
	e.POST("/deploy/clean/namespace", clean_deploy_request.TopologyCleanNamespaceHandler)
	e.POST("/deploy/replace", replace_topology.TopologyDeploymentReplaceHandler)
	e.POST("/deploy/status", deployment_status.TopologyDeploymentStatusHandler)
	e.POST("/deploy/cluster/status", deployment_status.ClusterDeploymentStatusHandler)

	// DELETE
	e.POST("/deploy/destroy", destroy_deploy_request.TopologyDestroyDeploymentHandler)
	e.POST("/deploy/ui/destroy", destroy_deploy_request.DestroyNamespaceHandler)
	e.POST("/resources/delete", destroy_deploy_request.DestroyResourceHandler)
	return e
}

func InternalRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	e = InternalDeployRoutes(e, k8Cfg)
	e = InternalDeployDestroyRoutes(e, k8Cfg)
	e = InternalDeployStatusRoutes(e)

	e = InternalSecretsRoutes(e, k8Cfg)
	return e
}

func InternalDeployStatusRoutes(e *echo.Group) *echo.Group {
	e.POST("/deploy/kns/create", workload_state.InsertOrUpdateWorkloadKnsStateHandler)
	e.POST("/deploy/kns/destroy", workload_state.DeleteWorkloadKnsStateHandler)
	e.POST("/deploy/status", workload_state.UpdateWorkloadStateHandler)
	return e
}

func InternalSecretsRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e.POST("/deploy/secrets", internal_secrets_deploy.DeploySecretsHandler)
	return e
}

func InternalDeployRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e.POST("/deploy/namespace", internal_deploy.DeployNamespaceHandler)
	e.POST("/deploy/deployment", internal_deploy.DeployDeploymentHandler)
	e.POST("/deploy/statefulset", internal_deploy.DeployStatefulSetHandler)
	e.POST("/deploy/configmap", internal_deploy.DeployConfigMapHandler)
	e.POST("/deploy/service", internal_deploy.DeployServiceHandler)
	e.POST("/deploy/ingress", internal_deploy.DeployIngressHandler)
	e.POST("/deploy/choreography/secrets", internal_deploy.DeployChoreographySecretsHandler)
	e.POST("/deploy/servicemonitor", internal_deploy.DeployServiceMonitorHandler)
	return e
}

func InternalDeployDestroyRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg
	e.POST("/deploy/destroy/namespace", internal_destroy_deploy.DestroyDeployNamespaceHandler)
	e.POST("/deploy/destroy/deployment", internal_destroy_deploy.DestroyDeployDeploymentHandler)
	e.POST("/deploy/destroy/statefulset", internal_destroy_deploy.DestroyDeployStatefulSetHandler)
	e.POST("/deploy/destroy/configmap", internal_destroy_deploy.DestroyDeployConfigMapHandler)
	e.POST("/deploy/destroy/service", internal_destroy_deploy.DestroyDeployServiceHandler)
	e.POST("/deploy/destroy/ingress", internal_destroy_deploy.DestroyDeployIngressHandler)
	e.POST("/deploy/destroy/servicemonitor", internal_destroy_deploy.DestroyDeployServiceMonitorHandler)
	return e
}
