package ethereum_slashing_protection_watermarking

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
)

func WatermarkAttestation(ctx context.Context, pubkey string, sourceEpoch, targetEpoch int) error {
	if IsSourceEpochGreaterThanTargetEpoch(ctx, pubkey, sourceEpoch, targetEpoch) {
		return errors.New("sourceEpoch greater than targetEpoch")
	}
	// Get Last Recorded Source/Target Epochs TODO
	prevSourceEpoch, prevTargetEpoch, err := FetchLastAttestation(ctx, pubkey)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Interface("pubkey", pubkey).Msg("failed to fetch last attestation")
		return err
	}
	if IsSourceEpochLessThanOrEqualToAnyPreviousAttestations(ctx, pubkey, sourceEpoch, prevSourceEpoch) {
		return errors.New("sourceEpoch less than or equal to a previous attestation sourceEpoch")
	}
	if IsTargetEpochLessThanOrEqualToAnyPreviousAttestations(ctx, pubkey, targetEpoch, prevTargetEpoch) {
		return errors.New("targetEpoch less than or equal to a previous attestation targetEpoch")
	}
	return nil
}

func FetchLastAttestation(ctx context.Context, pubkey string) (sourceEpoch, targetEpoch int, err error) {
	// TODO
	return 0, 0, nil
}

// TODO far future signing protection

// IsSourceEpochGreaterThanTargetEpoch if this is true then reject the signing request
func IsSourceEpochGreaterThanTargetEpoch(ctx context.Context, pubkey string, sourceEpoch, targetEpoch int) bool {
	if sourceEpoch >= targetEpoch {
		log.Ctx(ctx).Warn().Msgf("detected sourceEpoch %d greater than or equal targetEpoch %d for %s", sourceEpoch, targetEpoch, pubkey)
		return true
	}
	return false
}

// IsSourceEpochLessThanOrEqualToAnyPreviousAttestations if this is true then reject the signing request. hasSourceOlderThanWatermark
func IsSourceEpochLessThanOrEqualToAnyPreviousAttestations(ctx context.Context, pubkey string, newSourceEpoch, prevSourceEpoch int) bool {
	if newSourceEpoch <= prevSourceEpoch {
		log.Ctx(ctx).Warn().Msgf("detected newSourceEpoch %d less than or equal to maxRecordedSourceEpoch %d for %s", newSourceEpoch, prevSourceEpoch, pubkey)
		return true
	}
	return false
}

func IsTargetEpochLessThanOrEqualToAnyPreviousAttestations(ctx context.Context, pubkey string, newTargetEpoch, prevTargetEpoch int) bool {
	if newTargetEpoch <= prevTargetEpoch {
		log.Ctx(ctx).Warn().Msgf("detected newSourceEpoch %d less than or equal to maxRecordedSourceEpoch %d for %s", newTargetEpoch, prevTargetEpoch, pubkey)
		return true
	}
	return false
}
