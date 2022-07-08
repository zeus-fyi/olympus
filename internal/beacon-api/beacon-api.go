package beacon_api

import (
	"github.com/zeus-fyi/olympus/pkg/client"
	"github.com/zeus-fyi/olympus/pkg/utils"
)

var c client.Client

func init() {
	c.EnableBytesStrDecode = false
}

const getBlockByID = "eth/v2/beacon/blocks"
const getValidatorsByState = "eth/v1/beacon/states"

func GetValidatorsByState(beaconNode, stateID string) client.Reply {
	url := utils.UrlPathStrBuilder(beaconNode, getValidatorsByState, stateID, "validators")
	return c.Get(url)
}

func GetValidatorsBalancesByStateFilter(beaconNode, stateID string, valIndexes ...string) client.Reply {
	url := utils.UrlPathStrBuilder(beaconNode, getValidatorsByState, stateID, "validator_balances?id=")
	url = utils.UrlEncodeQueryParamList(url, valIndexes...)
	return c.Get(url)
}

func GetBlockByID(beaconNode, blockID string) client.Reply {
	url := utils.UrlPathStrBuilder(beaconNode, getBlockByID, blockID)
	return c.Get(url)
}
