package internal_deploy

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/internal/base_request"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func DeployDeploymentHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if request.Deployment != nil {
		_, err := zeus.K8Util.CreateDeploymentIfVersionLabelChangesOrDoesNotExist(ctx, request.Kns, request.Deployment, nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	} else {
		err := errors.New("no deployment workload was supplied")
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, nil)
}
