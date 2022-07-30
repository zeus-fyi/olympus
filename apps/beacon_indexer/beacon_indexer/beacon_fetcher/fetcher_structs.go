package beacon_fetcher

import (
	"github.com/zeus-fyi/olympus/pkg/beacon_api"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres/beacon_indexer/beacon_models"
)

type BeaconFetcher struct {
	NodeEndpoint string

	BeaconStateResults   beacon_api.ValidatorsStateBeacon
	BeaconBalanceResults beacon_api.ValidatorBalances

	Validators beacon_models.Validators
	Balances   beacon_models.ValidatorBalancesEpoch
}
