package internal_deploy

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func DeployServiceMonitorHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if request.Kns.TopologyBaseInfraWorkload.ServiceMonitor != nil {
		log.Debug().Interface("kns", request.Kns).Msg("DeployServiceMonitorHandler: CreateServiceMonitorIfVersionLabelChangesOrDoesNotExist")
		_, err := zeus.K8Util.CreateServiceMonitorIfVersionLabelChangesOrDoesNotExist(ctx, request.Kns.CloudCtxNs, request.Kns.TopologyBaseInfraWorkload.ServiceMonitor, nil)
		if err != nil {
			log.Err(err).Msg("DeployServiceMonitorHandler")
			return c.JSON(http.StatusInternalServerError, err)
		}
	} else {
		err := errors.New("no servicemonitor workload was supplied")
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, nil)
}
