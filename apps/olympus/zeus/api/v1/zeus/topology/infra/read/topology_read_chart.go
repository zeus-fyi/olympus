package read_infra

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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
	ou, ok := c.Get("orgUser").(org_users.OrgUser)
	if !ok {
		log.Error().Msg("ReadTopologyChart: orgUser not found")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	tr.OrgID = ou.OrgID
	tr.UserID = ou.UserID
	ctx := context.Background()
	err := tr.SelectTopology(ctx)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Msg("ReadTopologyChart: SelectTopology")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	nk := tr.GetTopologyBaseInfraWorkload()
	return c.JSON(http.StatusOK, nk)
}
