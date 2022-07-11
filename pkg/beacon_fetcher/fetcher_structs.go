package beacon_fetcher

import (
	beacon_models "github.com/zeus-fyi/olympus/databases/postgres/beacon-indexer/beacon-models"
	"github.com/zeus-fyi/olympus/internal/beacon-api/api_types"
)

type BeaconFetcher struct {
	NodeEndpoint string

	BeaconStateResults   api_types.ValidatorsStateBeacon
	BeaconBalanceResults api_types.ValidatorBalances

	Validators beacon_models.Validators
	Balances   beacon_models.ValidatorBalancesEpoch
}
