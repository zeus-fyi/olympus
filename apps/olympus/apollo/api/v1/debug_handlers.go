package v1

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/admin"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/beacon_indexer/beacon_models"
)

func DebugRequestHandler(c echo.Context) (err error) {
	log.Info().Msg("DebugRequestHandler")
	ctx := context.Background()
	var debug DebugReader
	//debug.ValidatorBalanceEntriesCount, err = beacon_models.SelectCountValidatorEpochBalanceEntries(ctx)
	//if err != nil {
	//	return c.JSON(http.StatusInternalServerError, "SelectCountValidatorEpochBalanceEntries had an error")
	//}
	//debug.ValidatorCount, err = beacon_models.SelectCountValidatorEntries(ctx)
	//if err != nil {
	//	return c.JSON(http.StatusInternalServerError, "SelectCountValidatorEntries had an error")
	//}

	checkpointEpoch := 164000
	chkPoint := beacon_models.ValidatorsEpochCheckpoint{}
	err = chkPoint.GetsOrderedNextEpochCheckpointWithBalancesRemainingAfterEpoch(ctx, checkpointEpoch)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "GetsOrderedNextEpochCheckpointWithBalancesRemainingAfterEpoch had an error")
	}
	debug.ForwardCheckpointEpoch = chkPoint.Epoch
	return c.JSON(http.StatusOK, debug)
}

func DebugGetPgConfigHandler(c echo.Context) (err error) {
	log.Info().Msg("DebugGetPgConfigHandler")
	ctx := context.Background()
	cfg := admin.ReadCfg(ctx)
	return c.JSON(http.StatusOK, cfg)
}

func DebugUpdatePgConfigHandler(c echo.Context) (err error) {
	log.Info().Msg("DebugUpdatePgConfigHandler")
	request := new(AdminDBConfigRequest)
	if err = c.Bind(request); err != nil {
		return err
	}
	ctx := context.Background()
	err = admin.UpdateConfigPG(ctx, request.ConfigChangePG)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "error updating config")
	}
	return c.JSON(http.StatusOK, "updated db config")
}
