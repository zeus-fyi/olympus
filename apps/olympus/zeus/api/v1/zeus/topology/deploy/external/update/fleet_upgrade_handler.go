package deploy_updates

import (
	"github.com/labstack/echo/v4"
)

func FleetUpgradeRequestHandler(c echo.Context) error {
	request := new(FleetUpgradeRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.UpgradeFleet(c)
}

func FleetRolloutRequestHandler(c echo.Context) error {
	request := new(FleetRolloutRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.RolloutRestartFleet(c)
}
