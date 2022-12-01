package beacon_api

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/client"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

var c client.Client

func init() {
	c.EnableBytesStrDecode = false
}

const getBlockByID = "eth/v2/beacon/blocks"
const getValidatorsByState = "eth/v1/beacon/states"

func GetValidatorsByState(ctx context.Context, beaconNode, stateID, status string) client.Reply {
	log.Info().Msg("BeaconAPI: GetValidatorsByState")
	if len(status) > 0 {
		status = fmt.Sprintf("?status=%s", status)
	}
	url := string_utils.UrlPathStrBuilder(beaconNode, getValidatorsByState, stateID, "validators", status)
	log.Debug().Interface("BeaconAPI: url:", url)
	bearer := "Bearer bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"

	m := make(map[string]string)
	m[echo.HeaderAuthorization] = bearer
	return c.Get(ctx, url, m)
}

func GetValidatorsFinalized(ctx context.Context, beaconNode string) client.Reply {
	log.Info().Msg("BeaconAPI: GetValidatorsFinalized")
	url := string_utils.UrlPathStrBuilder(beaconNode, getValidatorsByState, "finalized/validators")
	log.Debug().Interface("BeaconAPI: url:", url)
	bearer := "Bearer bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	m := make(map[string]string)
	m[echo.HeaderAuthorization] = bearer
	return c.Get(ctx, url, m)
}

func GetValidatorsByStateFilter(ctx context.Context, beaconNode, stateID string, encodedQueryURL, status string) client.Reply {
	log.Info().Msg("BeaconAPI: GetValidatorsByStateFilter")
	if len(status) > 0 {
		status = fmt.Sprintf("&status=%s", status)
	}

	url := string_utils.UrlPathStrBuilder(beaconNode, getValidatorsByState, stateID, "validators?"+encodedQueryURL+status)
	log.Debug().Interface("BeaconAPI: url:", url)
	bearer := "Bearer bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	m := make(map[string]string)
	m[echo.HeaderAuthorization] = bearer
	return c.Get(ctx, url, m)
}

func GetAllValidatorBalancesByState(ctx context.Context, beaconNode, stateID string) client.Reply {
	log.Info().Msg("BeaconAPI: GetAllValidatorBalancesByState")
	url := string_utils.UrlPathStrBuilder(beaconNode, getValidatorsByState, stateID, "validator_balances")
	log.Debug().Interface("BeaconAPI: url:", url)
	bearer := "Bearer bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	m := make(map[string]string)
	m[echo.HeaderAuthorization] = bearer
	return c.Get(ctx, url, m)
}

func GetValidatorsBalancesByStateFilter(ctx context.Context, beaconNode, stateID string, encodedQueryURL string) client.Reply {
	log.Info().Msg("BeaconAPI: GetValidatorsBalancesByStateFilter")

	url := string_utils.UrlPathStrBuilder(beaconNode, getValidatorsByState, stateID, "validator_balances?"+encodedQueryURL)
	log.Debug().Interface("BeaconAPI: url:", url)
	bearer := "Bearer bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	m := make(map[string]string)
	m[echo.HeaderAuthorization] = bearer
	return c.Get(ctx, url, m)
}

func GetBlockByID(ctx context.Context, beaconNode, blockID string) client.Reply {
	url := string_utils.UrlPathStrBuilder(beaconNode, getBlockByID, blockID)
	bearer := "Bearer bEX2piPZkxUuKwSkqkLh4KghmA7ZNDQnB"
	m := make(map[string]string)
	m[echo.HeaderAuthorization] = bearer
	return c.Get(ctx, url, m)
}
