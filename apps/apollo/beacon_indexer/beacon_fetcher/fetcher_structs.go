package beacon_fetcher

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/beacon_indexer/beacon_models"
	"github.com/zeus-fyi/olympus/datastores/redis_app/beacon_indexer"
	"github.com/zeus-fyi/olympus/pkg/web3/apollo/ethereum/beacon_api"
)

type BeaconFetcher struct {
	NodeEndpoint string

	BeaconStateResults   beacon_api.ValidatorsStateBeacon
	BeaconBalanceResults beacon_api.ValidatorBalances

	Validators beacon_models.Validators
	Balances   beacon_models.ValidatorBalancesEpoch

	Cache beacon_indexer.FetcherCache
}
