package infra_routes

import (
	"github.com/labstack/echo/v4"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	create_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create"
	read_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/read"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func Routes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	e.POST("/infra/create", create_infra.CreateTopologyInfraActionRequestHandler)
	e.POST("/infra/class/create", create_infra.CreateTopologyClassActionRequestHandler)
	e.POST("/infra/class/bases/create", create_infra.UpdateTopologyClassActionRequestHandler)
	e.POST("/infra/class/skeleton/bases/create", create_infra.CreateTopologySkeletonBasesActionRequestHandler)

	e.POST("/infra/read/chart", read_infra.ReadTopologyChartContentsHandler)
	e.GET("/infra/read/topologies", read_infra.ReadTopologiesMetadataRequestHandler)
	e.GET("/infra/read/org/topologies", read_infra.ReadTopologiesOrgCloudCtxNsHandler)
	return e
}

func UIRoutes(e *echo.Group, k8Cfg autok8s_core.K8Util) *echo.Group {
	zeus.K8Util = k8Cfg

	e.GET("/infra/ui/private/app/:id", read_infra.ReadOrgAppDetails)
	e.GET("/infra/ui/private/apps", read_infra.ReadOrgApps)
	e.POST("/infra/ui/create", create_infra.CreateTopologyInfraActionFromUIRequestHandler)
	e.POST("/infra/ui/preview/create", create_infra.PreviewCreateTopologyInfraActionRequestHandler)
	return e
}
