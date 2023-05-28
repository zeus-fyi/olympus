package deploy_updates

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type FleetUpgradeRequest struct {
	ClusterName string `json:"clusterName,omitempty"`
}

func (t *FleetUpgradeRequest) UpgradeFleet(c echo.Context) error {
	log.Debug().Msg("UpgradeFleet")

	return nil
}
