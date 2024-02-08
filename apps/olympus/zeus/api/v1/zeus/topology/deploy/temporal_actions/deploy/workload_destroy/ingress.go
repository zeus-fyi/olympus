package internal_destroy_deploy

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	autok8s_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
)

func DestroyDeployIngressHandlerWrapper(k autok8s_core.K8Util) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := context.Background()
		request := new(base_request.InternalDeploymentActionRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		if request.Kns.TopologyBaseInfraWorkload.Ingress != nil {
			err := k.DeleteIngressWithKns(ctx, request.Kns.CloudCtxNs, request.Kns.TopologyBaseInfraWorkload.Ingress.Name, nil)
			if err != nil {
				log.Err(err).Msg("DestroyDeployIngressHandler")
				return c.JSON(http.StatusInternalServerError, err)
			}
		} else {
			err := errors.New("no ingress workload was supplied")
			return c.JSON(http.StatusBadRequest, err)
		}
		return c.JSON(http.StatusOK, nil)
	}
}
