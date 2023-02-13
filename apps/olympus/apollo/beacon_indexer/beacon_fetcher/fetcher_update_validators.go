package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/beacon_indexer/beacon_models"
)

var (
	UpdateAllValidatorTimeout = time.Minute * 1
	UpdateValidatorBatchSize  = 1000
)

// UpdateAllValidators Routine TWO
func UpdateAllValidators() {
	log.Info().Msg("UpdateAllValidators")
	for {
		timeBegin := time.Now()
		err := fetchAllValidatorsToUpdate(context.Background(), UpdateAllValidatorTimeout)
		log.Err(err)
		log.Info().Interface("UpdateAllValidators took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(UpdateAllValidatorTimeout)
	}
}

func fetchAllValidatorsToUpdate(ctx context.Context, contextTimeout time.Duration) error {
	log.Info().Msg("fetchAllValidatorsToUpdate")
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	err := Fetcher.BeaconUpdateAllValidatorStates(ctxTimeout)
	log.Info().Err(err).Msg("UpdateAllValidators: fetchAllValidatorsToUpdate")
	return err
}

func (f *BeaconFetcher) BeaconUpdateAllValidatorStates(ctx context.Context) (err error) {
	log.Info().Msg("BeaconFetcher: BeaconUpdateAllValidatorStates")
	vsb, err := f.BeaconStateResults.FetchAllStateAndDecode(ctx, f.NodeEndpoint, "finalized", "")
	if err != nil {
		log.Err(err).Msg("BeaconUpdateValidatorStates: FetchStateAndDecode")
		return err
	}
	f.Validators = beacon_models.ToBeaconModelFormat(vsb)
	log.Info().Msg("BeaconFetcher: ToBeaconModelFormat")
	rowsUpdated, err := f.Validators.UpdateValidatorsFromBeaconAPI(ctx)
	log.Info().Msgf("BeaconFetcher: UpdateValidatorsFromBeaconAPI updated %d validators", rowsUpdated)
	if err != nil {
		log.Err(err).Msg("BeaconFetcher: UpdateValidatorsFromBeaconAPI")
		return err
	}
	if rowsUpdated <= 0 {
		log.Info().Msg("No validators were update")
	}
	return err
}

// FetchBeaconUpdateValidatorStates Routine ONE
func FetchBeaconUpdateValidatorStates(networkID int) {
	log.Info().Msg("FetchBeaconUpdateValidatorStates")

	for {
		timeBegin := time.Now()
		err := fetchValidatorsToUpdate(context.Background(), UpdateValidatorBatchSize, networkID, UpdateAllValidatorTimeout)
		log.Err(err)
		log.Info().Interface("FetchBeaconUpdateValidatorStates took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(1 * time.Second)
	}
}

func fetchValidatorsToUpdate(ctx context.Context, batchSize, networkID int, contextTimeout time.Duration) error {
	log.Info().Msg("fetchValidatorsToUpdate")

	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()
	err := Fetcher.BeaconUpdateValidatorStates(ctxTimeout, batchSize, networkID)
	log.Info().Err(err).Msg("fetchValidatorsToUpdate")
	return err
}

func (f *BeaconFetcher) BeaconUpdateValidatorStates(ctx context.Context, batchSize, networkID int) (err error) {
	log.Info().Msgf("BeaconFetcher: BeaconUpdateValidatorStates, batch size %d", batchSize)

	log.Info().Msg("BeaconUpdateValidatorStates: SelectValidatorsQueryOngoingStatesIndexesURLEncoded")
	indexes, err := beacon_models.SelectValidatorsQueryOngoingStatesIndexesURLEncoded(ctx, batchSize, networkID)
	if err != nil {
		log.Err(err).Msg("BeaconUpdateValidatorStates: SelectValidatorsQueryOngoingStatesIndexesURLEncoded")
		return err
	}

	indexCount := len(indexes)
	if indexCount <= 0 {
		log.Info().Msg("BeaconUpdateValidatorStates: had no new indexes")
		return nil
	}

	log.Info().Msg("BeaconUpdateValidatorStates: FetchStateAndDecode")
	vsb, err := f.BeaconStateResults.FetchStateAndDecode(ctx, f.NodeEndpoint, "finalized", indexes, "")
	if err != nil {
		log.Err(err).Msg("BeaconUpdateValidatorStates: FetchStateAndDecode")
		return err
	}
	f.Validators = beacon_models.ToBeaconModelFormat(vsb)
	rowsUpdated, err := f.Validators.UpdateValidatorsFromBeaconAPI(ctx)
	log.Info().Msgf("BeaconFetcher: UpdateValidatorsFromBeaconAPI updated %d validators", rowsUpdated)
	if err != nil {
		log.Err(err).Msg("BeaconFetcher: UpdateValidatorsFromBeaconAPI")
		return err
	}
	if rowsUpdated <= 0 {
		log.Info().Msg("No validators were update")
	}
	return err
}
