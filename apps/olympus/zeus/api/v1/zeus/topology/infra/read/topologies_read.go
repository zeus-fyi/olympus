package read_infra

import (
	"context"

	"github.com/labstack/echo/v4"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	base_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/base"
)

type TopologyActionReadRequest struct {
	base_infra.TopologyInfraActionRequest
	TopologyID int
}

func (t *TopologyActionReadRequest) ReadTopology(c echo.Context) error {
	tr := read_topology.NewInfraTopologyReader()
	tr.TopologyID = t.TopologyID
	tr.OrgID = t.OrgID
	tr.UserID = t.UserID

	ctx := context.Background()
	err := tr.SelectTopology(ctx)
	if err != nil {
		return err
	}
	return err
}
