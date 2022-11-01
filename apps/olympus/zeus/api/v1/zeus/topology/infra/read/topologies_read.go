package read_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/base"
)

type TopologyActionReadRequest struct {
	base.TopologyActionRequest
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
		return c.JSON(http.StatusInternalServerError, err)
	}
	nk := tr.GetNativeK8s()
	return c.JSON(http.StatusOK, nk)
}
