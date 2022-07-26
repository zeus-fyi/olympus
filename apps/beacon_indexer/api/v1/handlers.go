package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_fetcher"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
	beacon_models "github.com/zeus-fyi/olympus/pkg/datastores/postgres/beacon_indexer/beacon_models"
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

func HandleAdminGetRequest(c echo.Context) error {
	log.Info().Msg("HandleAdminGetRequest")
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

type DebugReader struct {
	ValidatorCount               int
	ValidatorBalanceEntriesCount int
}

func HandleDebugRequest(c echo.Context) (err error) {
	log.Info().Msg("HandleDebugRequest")
	ctx := context.Background()
	var debug DebugReader
	debug.ValidatorBalanceEntriesCount, err = beacon_models.SelectCountValidatorEpochBalanceEntries(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "SelectCountValidatorEpochBalanceEntries had an error")
	}
	debug.ValidatorCount, err = beacon_models.SelectCountValidatorEntries(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "SelectCountValidatorEntries had an error")
	}
	return c.JSON(http.StatusOK, debug)
}

func HandleDebugGetPgConfig(c echo.Context) (err error) {
	log.Info().Msg("HandleDebugGetPgConfig")
	ctx := context.Background()
	cfg := postgres.ReadCfg(ctx)
	return c.JSON(http.StatusOK, cfg)
}

func HandleDebugUpdatePgConfig(c echo.Context) (err error) {
	log.Info().Msg("HandleDebugUpdatePgConfig")
	request := new(AdminDBConfigRequest)
	if err = c.Bind(request); err != nil {
		return err
	}
	ctx := context.Background()
	err = postgres.UpdateConfigPG(ctx, request.ConfigChangePG)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error updating config")
	}
	return c.JSON(http.StatusOK, "updated db config")
}

func HandleDebugPgStats(c echo.Context) (err error) {
	log.Info().Msg("HandleDebugPgStats")
	ctx := context.Background()
	stats := postgres.Pg.PoolStats(ctx)
	return c.JSON(http.StatusOK, stats)
}

func HandlePingDB(c echo.Context) (err error) {
	log.Info().Msg("HandlePingDB")
	ctx := context.Background()
	err = postgres.Pg.Ping(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "pinged db successfully")
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
