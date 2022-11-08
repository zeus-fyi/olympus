package workload_state

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	create_topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/state"
)

// UpdateWorkloadStateHandler TODO must verify this is auth is scoped to user only
func UpdateWorkloadStateHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(InternalWorkloadStatusUpdateRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	status := create_topology_deployment_status.NewCreateState()
	status.TopologyStatus = request.TopologyStatus
	status.TopologyID = request.TopologyID
	err := status.InsertStatus(ctx)
	if err != nil {
		log.Err(err).Msg("UpdateWorkloadStateHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}
