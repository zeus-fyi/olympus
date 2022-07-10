package beacon_fetcher

import beacon_models "github.com/zeus-fyi/olympus/databases/postgres/beacon-indexer/beacon-models"

type BeaconFetcher struct {
	Validators beacon_models.Validators
}
