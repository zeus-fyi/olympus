package ethereum_slashing_protection_watermarking

import (
	"context"
)

func WatermarkAttestation(ctx context.Context, pubkey string, sourceEpoch, targetEpoch int) error {
	// TODO
	IsSourceEpochGreaterThanTargetEpoch(pubkey, sourceEpoch, targetEpoch)
	return nil
}
