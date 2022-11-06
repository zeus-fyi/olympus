package internal_deploy

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/internal/base_request"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func DeployNamespaceHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	_, err := zeus.K8Util.CreateNamespaceIfDoesNotExist(ctx, request.Kns)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
