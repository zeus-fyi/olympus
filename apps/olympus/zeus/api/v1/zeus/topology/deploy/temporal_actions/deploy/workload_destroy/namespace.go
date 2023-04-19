package internal_destroy_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func DestroyDeployNamespaceHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	log.Ctx(ctx).Info().Interface("request", request).Msg("DestroyDeployNamespaceHandler")
	err := zeus.K8Util.DeleteNamespace(ctx, request.Kns.CloudCtxNs)
	if err != nil {
		log.Err(err).Msg("DestroyDeployNamespaceHandler")
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
