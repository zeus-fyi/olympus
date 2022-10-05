package v1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_fetcher"
	"github.com/zeus-fyi/olympus/datastores/postgres"
	"github.com/zeus-fyi/olympus/pkg/logging"
)

type AdminConfigRequest struct {
	AdminConfig
}

type AdminDBConfigRequest struct {
	postgres.ConfigChangePG
}

type AdminConfig struct {
	LogLevel *zerolog.Level

	ValidatorBatchSize         *int
	ValidatorBalancesBatchSize *int
	ValidatorBalancesTimeout   *time.Duration
}

type AdminConfigReader struct {
	LogLevel zerolog.Level

	ValidatorBatchSize         int
	ValidatorBalancesBatchSize int
	ValidatorBalancesTimeout   time.Duration
}

func HandleAdminConfigRequest(c echo.Context) error {
	log.Info().Msg("HandleAdminConfigRequest")
	request := new(AdminConfigRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	if request.LogLevel != nil {
		ll := *request.LogLevel
		logging.SetLoggerLevel(ll)
		globalLevel := zerolog.GlobalLevel()

		if globalLevel != ll {
			log.Info().Msgf("level requested %s, level actual %s", ll, globalLevel)
			return c.JSON(http.StatusInternalServerError, "level did not update global logging level")
		}
		log.Info().Msgf("Set logging level to : %s", globalLevel)
	}

	if request.ValidatorBatchSize != nil {
		batchSize := *request.ValidatorBatchSize
		beacon_fetcher.NewValidatorBatchSize = batchSize
		log.Info().Msgf("Set ValidatorBatchSize level to : %d", batchSize)
	}

	if request.ValidatorBalancesBatchSize != nil {
		batchSize := *request.ValidatorBalancesBatchSize
		beacon_fetcher.NewValidatorBalancesBatchSize = batchSize
		log.Info().Msgf("Set NewValidatorBalancesBatchSize level to : %d", batchSize)
	}

	if request.ValidatorBalancesTimeout != nil {
		timeOut := *request.ValidatorBalancesTimeout
		beacon_fetcher.NewValidatorBalancesTimeout = timeOut
		log.Info().Msgf("Set ValidatorBalancesTimeout level to : %s", timeOut)
	}
	return c.JSON(http.StatusOK, "successfully updated config values")
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
