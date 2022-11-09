package read_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
)

type TopologyReadRequest struct {
	TopologyID int `json:"topologyID"`
}

func (t *TopologyReadRequest) ReadTopologyChart(c echo.Context) error {
	tr := read_topology.NewInfraTopologyReader()
	tr.TopologyID = t.TopologyID
	// from auth lookup
	ou := c.Get("orgUser").(org_users.OrgUser)
	tr.OrgID = ou.OrgID
	tr.UserID = ou.UserID
	ctx := context.Background()
	err := tr.SelectTopology(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	nk := tr.GetNativeK8s()
	return c.JSON(http.StatusOK, nk)
}
