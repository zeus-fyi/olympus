package beacon_api

import (
	"fmt"

	"bitbucket.org/zeus/eth-indexer/pkg/client"
)

var c client.Client

const getBlockByID = "/eth/v2/beacon/blocks/"
const getValidatorsByState = "/eth/v1/beacon/states"

func GetData(beaconNode, endpoint string) client.Reply {
	url := fmt.Sprintf("%s%s", beaconNode, endpoint)
	return c.Get(url)
}

func GetValidatorsByState(beaconNode, stateID string) client.Reply {
	url := fmt.Sprintf("%s%s/%s/validators", beaconNode, getValidatorsByState, stateID)
	return c.Get(url)
}

func GetBlockByID(beaconNode, blockID string) client.Reply {
	url := fmt.Sprintf("%s%s%s", beaconNode, getBlockByID, blockID)
	return c.Get(url)
}
