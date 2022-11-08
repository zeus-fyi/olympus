package destoy_deploy

import (
	"github.com/labstack/echo/v4"
)

func TopologyDestroyDeploymentHandler(c echo.Context) error {
	request := new(TopologyDestroyDeployRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.DestroyDeployedTopology(c)
}
