package internal_deploy

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
)

func DeployServiceHandlerWrapper(k autok8s_core.K8Util) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := context.Background()
		k8CfgInterface := c.Get("k8Cfg")
		if k8CfgInterface == nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Kubernetes configuration not found in context"})
		}
		k8Cfg, ok := k8CfgInterface.(autok8s_core.K8Util) // Ensure the type assertion is correct
		if ok {
			k = k8Cfg
		}
		// Attempt to retrieve the InternalDeploymentActionRequest from the context
		requestInterface := c.Get("internalDeploymentActionRequest")
		if requestInterface == nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal deployment action request not found in context"})
		}
		request, ok := requestInterface.(*base_request.InternalDeploymentActionRequest) // Ensure the type assertion is correct
		if !ok {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid request type"})
		}
		if request.Kns.TopologyBaseInfraWorkload.Service != nil {
			log.Debug().Interface("kns", request.Kns).Msg("DeployServiceHandler: CreateServiceIfVersionLabelChangesOrDoesNotExist")
			_, err := k.CreateServiceIfVersionLabelChangesOrDoesNotExist(ctx, request.Kns.CloudCtxNs, request.Kns.TopologyBaseInfraWorkload.Service, nil)
			if err != nil {
				log.Err(err).Msg("DeployServiceHandler")
				return c.JSON(http.StatusInternalServerError, err)
			}
		} else {
			err := errors.New("no service workload was supplied")
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, nil)
	}
}
