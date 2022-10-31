package create_infra

import (
	"github.com/labstack/echo/v4"
	base_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/base"
)

type TopologyActionCreateRequest struct {
	base_infra.TopologyInfraActionRequest
}

func (t *TopologyActionCreateRequest) CreateTopology(c echo.Context) error {

	//ctx := context.Background()
	//q := sql_query_templates.NewQueryParam("InsertTopology", "topologies", "where", 1000, []string{})
	//top := topology.NewInfrastructureTopologyClass()

	return nil
}
