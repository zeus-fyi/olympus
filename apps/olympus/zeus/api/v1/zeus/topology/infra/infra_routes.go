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

	e.POST("/infra/read/chart", read_infra.ReadTopologyChartContentsHandler)
	e.GET("/infra/read/topologies", read_infra.ReadTopologiesMetadataRequestHandler)
	return e
}
