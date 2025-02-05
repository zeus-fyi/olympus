package deployment_status

import (
	"github.com/labstack/echo/v4"
)

func TopologyDeploymentStatusHandler(c echo.Context) error {
	request := new(TopologyDeploymentStatusRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadDeployedTopologyStatuses(c)
}

func ClusterDeploymentStatusHandler(c echo.Context) error {
	request := new(ClusterDeploymentStatusRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.ReadDeployedClusterTopologies(c)
}
