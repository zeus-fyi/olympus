package workload_state

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	create_kns "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/kns"
)

// UpdateWorkloadKnsStateHandler TODO must verify this is auth is scoped to user only
func UpdateWorkloadKnsStateHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(kns.TopologyKubeCtxNs)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	err := create_kns.InsertKns(ctx, request)
	if err != nil {
		log.Err(err).Interface("kns", request).Msg("UpdateWorkloadKnsStateHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, request)
}
