package internal_deploy

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/internal/base_request"
	"github.com/zeus-fyi/olympus/zeus/pkg/zeus"
)

func DeployConfigMapHandler(c echo.Context) error {
	ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if request.ConfigMap != nil {
		_, err := zeus.K8Util.CreateConfigMapIfVersionLabelChangesOrDoesNotExist(ctx, request.Kns, request.ConfigMap, nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	} else {
		err := errors.New("no configmap workload was supplied")
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, nil)
}
