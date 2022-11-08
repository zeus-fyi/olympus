package read_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
)

type TopologyReadRequest struct {
	TopologyID int `json:"topologyID"`
}

func (t *TopologyReadRequest) ReadTopology(c echo.Context) error {
	tr := read_topology.NewInfraTopologyReader()
	tr.TopologyID = t.TopologyID
	// from auth lookup
	orgID := c.Get("orgID")
	tr.OrgID = orgID.(int)

	userID := c.Get("userID")
	tr.UserID = userID.(int)

	ctx := context.Background()
	err := tr.SelectTopology(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	nk := tr.GetNativeK8s()
	return c.JSON(http.StatusOK, nk)
}
