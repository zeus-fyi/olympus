package ethereum_slashing_protection_watermarking

import (
	"github.com/rs/zerolog/log"
)

// IsSourceEpochGreaterThanTargetEpoch if this is true then reject the signing request
func IsSourceEpochGreaterThanTargetEpoch(pubkey string, sourceEpoch, targetEpoch int) bool {
	if sourceEpoch > targetEpoch {
		log.Warn().Msgf("detected sourceEpoch %d greater than targetEpoch %d for %s", sourceEpoch, targetEpoch, "pubkey")
		return true
	}
	return false
}

func RecordAttestationSignatureAtEpoch(pubkey string, epoch int) bool {
	return false
}

func RecordBlockSignatureAtSlot(pubkey string, source, target int) bool {
	if source > target {
		log.Warn().Msgf("detected sourceEpoch %d greater than targetEpoch %d for %s", source, target, "pubkey")
		return true
	}
	return false
}
