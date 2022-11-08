package destoy_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
)

type TopologyDestroyDeployRequest struct {
	TopologyID int
}

func (t *TopologyDestroyDeployRequest) DestroyDeployedTopology(c echo.Context) error {
	tr := read_topology.NewInfraTopologyReader()
	tr.TopologyID = t.TopologyID
	ctx := context.Background()
	err := tr.SelectTopology(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return err
}
