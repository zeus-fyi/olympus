package beacon_fetcher

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/beacon_indexer/beacon_models"
	"github.com/zeus-fyi/olympus/datastores/redis/apps/beacon_indexer"
	"github.com/zeus-fyi/olympus/pkg/apollo/ethereum/consensus_client_apis/beacon_api"
)

type BeaconFetcher struct {
	NodeEndpoint string

	BeaconStateResults   beacon_api.ValidatorsStateBeacon
	BeaconBalanceResults beacon_api.ValidatorBalances

	Validators beacon_models.Validators
	Balances   beacon_models.ValidatorBalancesEpoch

	Cache beacon_indexer.FetcherCache
}
