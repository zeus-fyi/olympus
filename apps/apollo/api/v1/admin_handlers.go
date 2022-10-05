package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_fetcher"
)

type DebugReader struct {
	ValidatorCount               int
	ValidatorBalanceEntriesCount int
	ForwardCheckpointEpoch       int
}

func AdminGetRequestHandler(c echo.Context) error {
	log.Info().Msg("AdminGetRequestHandler")
	request := new(AdminConfigRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	cfgRead := AdminConfigReader{
		LogLevel:                   zerolog.GlobalLevel(),
		ValidatorBatchSize:         beacon_fetcher.NewValidatorBatchSize,
		ValidatorBalancesBatchSize: beacon_fetcher.NewValidatorBalancesBatchSize,
		ValidatorBalancesTimeout:   beacon_fetcher.NewValidatorBalancesTimeout,
	}
	return c.JSON(http.StatusOK, cfgRead)
}
