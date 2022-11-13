package workload_state

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	create_topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/state"
)

// UpdateWorkloadStateHandler TODO must verify this is auth is scoped to user only
func UpdateWorkloadStateHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(topology_deployment_status.Status)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	status := topology_deployment_status.NewTopologyStatus()
	status.TopologyStatus = request.TopologyStatus
	status.TopologyID = request.TopologyID
	err := create_topology_deployment_status.InsertOrUpdateStatus(ctx, &status)
	if err != nil {
		log.Err(err).Interface("status", status).Msg("UpdateWorkloadStateHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, status)
}
