package internal_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func DeployJobHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	log.Debug().Interface("kns", request.Kns).Msg("DeployJobHandler")
	_, err := zeus.K8Util.CreateJob(ctx, request.Kns.CloudCtxNs, request.Kns.TopologyBaseInfraWorkload.Job)
	if err != nil {
		log.Err(err).Msg("DeployJobHandler")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
