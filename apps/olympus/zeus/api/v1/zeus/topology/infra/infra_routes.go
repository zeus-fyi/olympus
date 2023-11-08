package infra_routes

import (
	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	create_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create"
	read_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/read"
	update_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/update"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func Routes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	e.POST("/infra/create", create_infra.CreateTopologyInfraActionRequestHandler)
	e.POST("/infra/class/create", create_infra.CreateTopologyClassActionRequestHandler)
	e.POST("/infra/class/bases/create", create_infra.UpdateTopologyClassActionRequestHandler)
	e.POST("/infra/class/skeleton/bases/create", create_infra.CreateTopologySkeletonBasesActionRequestHandler)

	// matrix
	e.POST("/infra/matrix/create", create_infra.CreateMatrixInfraActionRequestHandler)

	e.POST("/infra/read/chart", read_infra.ReadTopologyChartContentsHandler)
	e.GET("/infra/read/topologies", read_infra.ReadTopologiesMetadataRequestHandler)
	e.GET("/infra/read/org/topologies", read_infra.ReadTopologiesOrgCloudCtxNsHandler)
	e.GET("/infra/read/org/topologies/apps", read_infra.ReadClusterAppViewOrgCloudCtxNsHandler)

	return e
}

func UIRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	// matrix
	e.GET("/infra/ui/matrix/public/apps/:name", read_infra.PublicAppsMatrixRequestHandler)

	e.GET("/infra/ui/apps/microservice", read_infra.MicroserviceAppsHandler)
	e.GET("/infra/ui/apps/avax", read_infra.AvaxAppsHandler)
	e.GET("/infra/ui/apps/sui", read_infra.SuiAppsHandler)
	e.GET("/infra/ui/apps/eth/beacon/ephemeral", read_infra.EthAppsHandler)
	e.GET("/infra/ui/private/app/:id", read_infra.ReadOrgAppDetailsHandler)
	e.GET("/infra/ui/private/apps", read_infra.ReadOrgAppsHandler)
	e.POST("/infra/ui/cluster/create", create_infra.CreateTopologyInfraActionFromUIRequestHandler)
	e.POST("/infra/ui/cluster/preview", create_infra.PreviewCreateTopologyInfraActionRequestHandler)
	e.POST("/infra/ui/cluster/update", update_infra.UpdateClusterInfraRequestHandler)
	return e
}
