package internal_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/base_request"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func DeployNamespaceHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	kns := zeus_core.NewKubeCtxNsFromTopologyKns(request.Kns)
	log.Debug().Interface("kns", kns).Msg("DeployNamespaceHandler: CreateNamespaceIfDoesNotExist")
	_, err := zeus.K8Util.CreateNamespaceIfDoesNotExist(ctx, kns)
	if err != nil {
		log.Err(err).Msg("DeployNamespaceHandler")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
