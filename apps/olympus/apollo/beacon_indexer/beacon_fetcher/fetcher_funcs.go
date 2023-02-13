package beacon_fetcher

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/beacon_indexer/beacon_models"
)

func (f *BeaconFetcher) BeaconFindNewAndMissingValidatorIndexes(ctx context.Context, batchSize, networkID int) (err error) {
	log.Info().Msg("BeaconFetcher: BeaconFindNewAndMissingValidatorIndexes")

	log.Info().Msg("BeaconFindNewAndMissingValidatorIndexes: FindNewValidatorsToQueryBeaconURLEncoded")
	indexes, err := beacon_models.FindNewValidatorsToQueryBeaconURLEncoded(ctx, batchSize, networkID)
	if err != nil {
		log.Err(err).Msg("BeaconStateUpdater: FindNewValidatorsToQueryBeaconURLEncoded")
		return err
	}
	if len(indexes) <= 0 {
		log.Info().Msg("FindNewValidatorsToQueryBeaconURLEncoded: had no new indexes")
		return nil
	}
	log.Info().Msg("BeaconFindNewAndMissingValidatorIndexes: FetchStateAndDecode")
	vsb, err := f.BeaconStateResults.FetchStateAndDecode(ctx, f.NodeEndpoint, "finalized", indexes, "")
	if err != nil {
		log.Err(err).Msg("BeaconFindNewAndMissingValidatorIndexes: FetchStateAndDecode")
		return err
	}
	f.Validators = beacon_models.ToBeaconModelFormat(vsb)
	log.Info().Msg("BeaconFindNewAndMissingValidatorIndexes: InsertValidatorsFromBeaconAPI")
	err = f.Validators.InsertValidatorsFromBeaconAPI(ctx)
	if err != nil {
		log.Err(err).Msg("BeaconFindNewAndMissingValidatorIndexes: InsertValidatorsFromBeaconAPI")
		return err
	}
	return err
}
