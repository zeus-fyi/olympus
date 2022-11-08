package create_or_update_deploy

import (
	"github.com/labstack/echo/v4"
)

func TopologyDeploymentHandler(c echo.Context) error {
	request := new(TopologyDeployRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.DeployTopology(c)
}
