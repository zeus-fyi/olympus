package internal_destroy_deploy

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/base_request"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func DestroyDeployIngressHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if request.Ingress != nil {
		kns := zeus_core.NewKubeCtxNsFromTopologyKns(request.Kns)
		err := zeus.K8Util.DeleteIngressWithKns(ctx, kns, request.Ingress.Name, nil)
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
