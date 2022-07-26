package beacon_api

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/client"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

var c client.Client

func init() {
	c.EnableBytesStrDecode = false
}

const getBlockByID = "eth/v2/beacon/blocks"
const getValidatorsByState = "eth/v1/beacon/states"

func GetValidatorsByState(ctx context.Context, beaconNode, stateID string) client.Reply {
	log.Info().Msg("BeaconAPI: GetValidatorsByState")
	url := string_utils.UrlPathStrBuilder(beaconNode, getValidatorsByState, stateID, "validators")
	log.Debug().Interface("BeaconAPI: url:", url)
	return c.Get(ctx, url)
}

func GetValidatorsByStateFilter(ctx context.Context, beaconNode, stateID string, encodedQueryURL string) client.Reply {
	log.Info().Msg("BeaconAPI: GetValidatorsByStateFilter")
	url := string_utils.UrlPathStrBuilder(beaconNode, getValidatorsByState, stateID, "validators?id="+encodedQueryURL)
	log.Debug().Interface("BeaconAPI: url:", url)
	return c.Get(ctx, url)
}

func GetValidatorsBalancesByStateFilter(ctx context.Context, beaconNode, stateID string, encodedQueryURL string) client.Reply {
	log.Info().Msg("BeaconAPI: GetValidatorsBalancesByStateFilter")
	url := string_utils.UrlPathStrBuilder(beaconNode, getValidatorsByState, stateID, "validator_balances?id="+encodedQueryURL)
	log.Debug().Interface("BeaconAPI: url:", url)
	return c.Get(ctx, url)
}

func GetBlockByID(ctx context.Context, beaconNode, blockID string) client.Reply {
	url := string_utils.UrlPathStrBuilder(beaconNode, getBlockByID, blockID)
	return c.Get(ctx, url)
}
