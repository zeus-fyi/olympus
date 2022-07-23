package beacon_fetcher

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_api/api_types"
	beacon_models2 "github.com/zeus-fyi/olympus/pkg/databases/postgres/beacon-indexer/beacon-models"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

func (f *BeaconFetcher) BeaconFindNewAndMissingValidatorIndexes(ctx context.Context, batchSize int) (err error) {
	log.Info().Msg("BeaconFetcher: BeaconFindNewAndMissingValidatorIndexes")

	log.Info().Msg("BeaconFindNewAndMissingValidatorIndexes: FindNewValidatorsToQueryBeaconURLEncoded")
	indexes, err := beacon_models2.FindNewValidatorsToQueryBeaconURLEncoded(ctx, batchSize)
	if err != nil {
		log.Error().Err(err).Msg("BeaconStateUpdater: FindNewValidatorsToQueryBeaconURLEncoded")
		return err
	}
	if len(indexes) <= 0 {
		log.Info().Msg("FindNewValidatorsToQueryBeaconURLEncoded: had no new indexes")
		return nil
	}

	log.Info().Msg("BeaconFindNewAndMissingValidatorIndexes: FetchStateAndDecode")
	err = f.BeaconStateResults.FetchStateAndDecode(ctx, f.NodeEndpoint, "finalized", indexes)
	if err != nil {
		log.Error().Err(err).Msg("BeaconFindNewAndMissingValidatorIndexes: FetchStateAndDecode")
		return err
	}
	f.Validators = beacon_models2.ToBeaconModelFormat(f.BeaconStateResults)
	log.Info().Msg("BeaconFindNewAndMissingValidatorIndexes: InsertValidatorsFromBeaconAPI")
	err = f.Validators.InsertValidatorsFromBeaconAPI(ctx)
	if err != nil {
		log.Error().Err(err).Msg("BeaconFindNewAndMissingValidatorIndexes: InsertValidatorsFromBeaconAPI")
		return err
	}
	return err
}

func (f *BeaconFetcher) BeaconUpdateValidatorStates(ctx context.Context, batchSize int) (err error) {
	log.Info().Msg("BeaconFetcher: BeaconUpdateValidatorStates")

	log.Info().Msg("BeaconUpdateValidatorStates: SelectValidatorsQueryOngoingStatesIndexesURLEncoded")
	indexes, err := beacon_models2.SelectValidatorsQueryOngoingStatesIndexesURLEncoded(ctx, batchSize)
	if err != nil {
		log.Error().Err(err).Msg("BeaconUpdateValidatorStates: SelectValidatorsQueryOngoingStatesIndexesURLEncoded")
		return err
	}
	if len(indexes) <= 0 {
		log.Info().Msg("BeaconUpdateValidatorStates: had no new indexes")
		return nil
	}

	log.Info().Msg("BeaconUpdateValidatorStates: FetchStateAndDecode")
	err = f.BeaconStateResults.FetchStateAndDecode(ctx, f.NodeEndpoint, "finalized", indexes)
	if err != nil {
		log.Error().Err(err).Msg("BeaconUpdateValidatorStates: FetchStateAndDecode")
		return err
	}
	f.Validators = beacon_models2.ToBeaconModelFormat(f.BeaconStateResults)
	log.Info().Msg("BeaconUpdateValidatorStates: UpdateValidatorsFromBeaconAPI")
	vals, err := f.Validators.UpdateValidatorsFromBeaconAPI(ctx)
	if err != nil {
		log.Error().Err(err).Msg("BeaconUpdateValidatorStates: InsertValidatorsFromBeaconAPI")
		return err
	}
	if len(vals.Validators) <= 0 {
		log.Info().Interface("No validators were returned", vals.Validators)
	}
	return err
}

func (f *BeaconFetcher) FindAndQueryAndUpdateValidatorBalances(ctx context.Context, batchSize int) error {
	log.Info().Msg("BeaconFetcher: FindAndQueryAndUpdateValidatorBalances")

	log.Info().Msg("FindAndQueryAndUpdateValidatorBalances: SelectValidatorsToQueryBeaconForBalanceUpdates")
	epochMap, err := beacon_models2.SelectValidatorsToQueryBalancesByEpoch(ctx, batchSize)
	if err != nil {
		log.Error().Err(err).Msg("FindAndQueryAndUpdateValidatorBalances: SelectValidatorsToQueryBeaconForBalanceUpdates")
		return err
	}

	for epoch, vbs := range epochMap {
		valBalances := beacon_models2.ValidatorBalancesEpoch{}
		valBalances.ValidatorBalance = vbs
		var beaconAPI api_types.ValidatorBalances
		slot := misc.ConvertEpochToSlot(epoch)
		log.Info().Interface("BeaconFetcher: Fetching Data at Slot:", slot)
		log.Info().Msg("BeaconFetcher: FetchStateAndDecode")
		err = beaconAPI.FetchStateAndDecode(ctx, f.NodeEndpoint, slot, valBalances.FormatValidatorBalancesEpochIndexesToURLList())
		if err != nil {
			log.Info().Interface("FormatValidatorBalancesEpochIndexesToURLList: ", valBalances.FormatValidatorBalancesEpochIndexesToURLList())
			log.Error().Err(err).Msg("FindAndQueryAndUpdateValidatorBalances: FetchStateAndDecode")
			return err
		}
		if len(beaconAPI.Data) <= 0 {
			log.Info().Interface("BeaconFetcher: FetchStateAndDecode returned zero balances for ", valBalances.FormatValidatorBalancesEpochIndexesToURLList())
			return nil
		}
		log.Info().Msg("BeaconFetcher: Convert API data to model format")
		valBalances = convertBeaconAPIBalancesToModelBalance(epoch, beaconAPI, valBalances)
		log.Info().Msg("BeaconFetcher: InsertValidatorBalancesForNextEpoch")
		err = valBalances.InsertValidatorBalancesForNextEpoch(ctx)
		if err != nil {
			log.Error().Err(err).Msg("FindAndQueryAndUpdateValidatorBalances: InsertValidatorBalancesForNextEpoch")
			return err
		}
	}
	return err
}

func convertBeaconAPIBalancesToModelBalance(epoch int64, beaconBalanceAPI api_types.ValidatorBalances, valBalances beacon_models2.ValidatorBalancesEpoch) beacon_models2.ValidatorBalancesEpoch {
	log.Info().Msg("BeaconFetcher: convertBeaconAPIBalancesToModelBalance")
	valBalances.ValidatorBalance = make([]beacon_models2.ValidatorBalanceEpoch, len(beaconBalanceAPI.Data))
	for i, beaconBalanceResult := range beaconBalanceAPI.Data {
		var epochResult beacon_models2.ValidatorBalanceEpoch
		epochResult.Epoch = epoch
		epochResult.Index = strings.Int64StringParser(beaconBalanceResult.Index)
		epochResult.TotalBalanceGwei = strings.Int64StringParser(beaconBalanceResult.Balance)
		valBalances.ValidatorBalance[i] = epochResult
	}
	return valBalances
}
