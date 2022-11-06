package workload_state

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/internal/base_request"
)

func UpdateWorkloadStateHandler(c echo.Context) error {
	//ctx := context.Background()
	request := new(base_request.InternalDeploymentActionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
