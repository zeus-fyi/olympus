package beacon_fetcher

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_api/api_types"
	beacon_models "github.com/zeus-fyi/olympus/pkg/datastores/postgres/beacon_indexer/beacon_models"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func (f *BeaconFetcher) BeaconFindNewAndMissingValidatorIndexes(ctx context.Context, batchSize int) (err error) {
	log.Info().Msg("BeaconFetcher: BeaconFindNewAndMissingValidatorIndexes")

	log.Info().Msg("BeaconFindNewAndMissingValidatorIndexes: FindNewValidatorsToQueryBeaconURLEncoded")
	indexes, err := beacon_models.FindNewValidatorsToQueryBeaconURLEncoded(ctx, batchSize)
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
	f.Validators = beacon_models.ToBeaconModelFormat(f.BeaconStateResults)
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
	indexes, err := beacon_models.SelectValidatorsQueryOngoingStatesIndexesURLEncoded(ctx, batchSize)
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
	rowsUpdated, err := f.Validators.UpdateValidatorsFromBeaconAPI(ctx)
	log.Info().Msgf("BeaconFetcher: UpdateValidatorsFromBeaconAPI updated %d validators", rowsUpdated)
	if err != nil {
		log.Error().Err(err).Msg("BeaconFetcher: UpdateValidatorsFromBeaconAPI")
		return err
	}
	if rowsUpdated <= 0 {
		log.Info().Msg("No validators were update")
	}
	return err
}

func (f *BeaconFetcher) FindAndQueryAndUpdateValidatorBalances(ctx context.Context, batchSize int) error {
	log.Info().Msg("BeaconFetcher: FindAndQueryAndUpdateValidatorBalances")

	log.Info().Msg("FindAndQueryAndUpdateValidatorBalances: SelectValidatorsToQueryBeaconForBalanceUpdates")
	nextEpochSlotMap, err := beacon_models.SelectValidatorsToQueryBalancesByEpochSlot(ctx, batchSize)
	if err != nil {
		log.Error().Err(err).Msg("FindAndQueryAndUpdateValidatorBalances: SelectValidatorsToQueryBeaconForBalanceUpdates")
		return err
	}

	for nextEpoch, vbs := range nextEpochSlotMap {
		valBalances := beacon_models.ValidatorBalancesEpoch{}
		valBalances.ValidatorBalances = vbs
		var beaconAPI api_types.ValidatorBalances
		nextEpochSlot := misc.ConvertEpochToSlot(nextEpoch)
		beaconAPI.Epoch = nextEpoch
		log.Info().Interface("BeaconFetcher: Fetching Data at Slot:", nextEpochSlot)
		log.Info().Msg("BeaconFetcher: FetchStateAndDecode")
		err = beaconAPI.FetchStateAndDecode(ctx, f.NodeEndpoint, nextEpochSlot, valBalances.FormatValidatorBalancesEpochIndexesToURLList())
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
		valBalances = convertBeaconAPIBalancesToModelBalance(beaconAPI, valBalances)
		log.Info().Msg("BeaconFetcher: InsertValidatorBalancesForNextEpoch")
		err = valBalances.InsertValidatorBalancesForNextEpoch(ctx)
		if err != nil {
			log.Error().Err(err).Msg("FindAndQueryAndUpdateValidatorBalances: InsertValidatorBalancesForNextEpoch")
			return err
		}
	}
	return err
}

func convertBeaconAPIBalancesToModelBalance(beaconBalanceAPI api_types.ValidatorBalances, valBalances beacon_models.ValidatorBalancesEpoch) beacon_models.ValidatorBalancesEpoch {
	log.Info().Msg("BeaconFetcher: convertBeaconAPIBalancesToModelBalance")
	valBalances.ValidatorBalances = make([]beacon_models.ValidatorBalanceEpoch, len(beaconBalanceAPI.Data))
	for i, beaconBalanceResult := range beaconBalanceAPI.Data {
		var epochResult beacon_models.ValidatorBalanceEpoch
		epochResult.Epoch = beaconBalanceAPI.Epoch
		epochResult.Index = string_utils.Int64StringParser(beaconBalanceResult.Index)
		epochResult.TotalBalanceGwei = string_utils.Int64StringParser(beaconBalanceResult.Balance)
		valBalances.ValidatorBalances[i] = epochResult
	}
	return valBalances
}

func (f *BeaconFetcher) FetchAllValidatorBalances(ctx context.Context, epoch int64) (beacon_models.ValidatorBalancesEpoch, error) {
	log.Info().Msg("BeaconFetcher: FetchAllValidatorBalancesAtSlot")
	var valBalances beacon_models.ValidatorBalancesEpoch
	var beaconAPI api_types.ValidatorBalances

	slot := epoch * 32
	slotToQuery := misc.ConvertEpochToSlot(slot)

	err := beaconAPI.FetchAllValidatorBalancesAtStateAndDecode(ctx, f.NodeEndpoint, slotToQuery)
	if err != nil {
		log.Error().Err(err).Msg("BeaconFetcher: QueryAllValidatorBalancesAtSlot")
		return valBalances, err
	}
	log.Info().Msg("BeaconFetcher: Convert API data to model format")
	valBalances.ValidatorBalances = make([]beacon_models.ValidatorBalanceEpoch, len(beaconAPI.Data))
	for i, vbFromAPI := range beaconAPI.Data {
		vbForDataEntry := beacon_models.ValidatorBalanceEpoch{
			Validator:        beacon_models.Validator{Index: string_utils.Int64StringParser(vbFromAPI.Index)},
			Epoch:            epoch,
			TotalBalanceGwei: string_utils.Int64StringParser(vbFromAPI.Balance),
		}
		valBalances.ValidatorBalances[i] = vbForDataEntry
	}

	return valBalances, nil
}

func (f *BeaconFetcher) BeaconUpdateAllValidatorStates(ctx context.Context) (err error) {
	log.Info().Msg("BeaconFetcher: BeaconUpdateAllValidatorStates")
	err = f.BeaconStateResults.FetchAllStateAndDecode(ctx, f.NodeEndpoint, "finalized")
	if err != nil {
		log.Error().Err(err).Msg("BeaconUpdateValidatorStates: FetchStateAndDecode")
		return err
	}
	f.Validators = beacon_models.ToBeaconModelFormat(f.BeaconStateResults)
	log.Info().Msg("BeaconFetcher: ToBeaconModelFormat")
	rowsUpdated, err := f.Validators.UpdateValidatorsFromBeaconAPI(ctx)
	log.Info().Msgf("BeaconFetcher: UpdateValidatorsFromBeaconAPI updated %d validators", rowsUpdated)
	if err != nil {
		log.Error().Err(err).Msg("BeaconFetcher: UpdateValidatorsFromBeaconAPI")
		return err
	}
	if rowsUpdated <= 0 {
		log.Info().Msg("No validators were update")
	}
	return err
}
