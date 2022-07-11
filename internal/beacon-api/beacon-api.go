package beacon_api

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/client"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

var c client.Client

func init() {
	c.EnableBytesStrDecode = false
}

const getBlockByID = "eth/v2/beacon/blocks"
const getValidatorsByState = "eth/v1/beacon/states"

func GetValidatorsByState(ctx context.Context, beaconNode, stateID string) client.Reply {
	url := strings.UrlPathStrBuilder(beaconNode, getValidatorsByState, stateID, "validators")
	return c.Get(ctx, url)
}

func GetValidatorsByStateFilter(ctx context.Context, beaconNode, stateID string, encodedQueryURL string) client.Reply {
	url := strings.UrlPathStrBuilder(beaconNode, getValidatorsByState, stateID, "validators?="+encodedQueryURL)
	return c.Get(ctx, url)
}

func GetValidatorsBalancesByStateFilter(ctx context.Context, beaconNode, stateID string, valIndexes ...string) client.Reply {
	url := strings.UrlPathStrBuilder(beaconNode, getValidatorsByState, stateID, "validator_balances?id=")
	url = strings.UrlEncodeQueryParamList(url, valIndexes...)
	return c.Get(ctx, url)
}

func GetBlockByID(ctx context.Context, beaconNode, blockID string) client.Reply {
	url := strings.UrlPathStrBuilder(beaconNode, getBlockByID, blockID)
	return c.Get(ctx, url)
}
