package beacon_fetcher

import (
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_api/api_types"
	beacon_models2 "github.com/zeus-fyi/olympus/pkg/datastores/postgres/beacon_indexer/beacon_models"
)

type BeaconFetcher struct {
	NodeEndpoint string

	BeaconStateResults   api_types.ValidatorsStateBeacon
	BeaconBalanceResults api_types.ValidatorBalances

	Validators beacon_models2.Validators
	Balances   beacon_models2.ValidatorBalancesEpoch
}
