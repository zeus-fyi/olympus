package workload_state

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	create_kns "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/kns"
	delete_kns "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/delete/topologies/topology/kns"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func InsertOrUpdateWorkloadKnsStateHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(zeus_req_types.TopologyDeployRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	err := create_kns.InsertKns(ctx, request)
	if err != nil {
		log.Err(err).Interface("kns", request).Msg("InsertOrUpdateWorkloadKnsStateHandler")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, request)
}

func DeleteWorkloadKnsStateHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(zeus_req_types.TopologyDeployRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if request.TopologyID == 0 {
		err := delete_kns.DeleteKnsByOrgAccessAndCloudCtx(ctx, request)
		if err != nil {
			log.Err(err).Interface("kns", request).Msg("DeleteWorkloadKnsStateHandler")
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else {
		err := delete_kns.DeleteKns(ctx, request)
		if err != nil {
			log.Err(err).Interface("kns", request).Msg("DeleteWorkloadKnsStateHandler")
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	return c.JSON(http.StatusOK, request)
}
