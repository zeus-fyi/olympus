package ethereum_slashing_protection_watermarking

import (
	"context"
	"strconv"

	"github.com/rs/zerolog/log"
	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
)

func RecordAttestationSignatureAtEpoch(pubkey string, source, target int) bool {
	return false
}

func RecordBlockSignatureAtSlot(pubkey string, slot int) bool {
	return false
}

func ConvertAttSourceTargetsToInt(ctx context.Context, att consensys_eth2_openapi.AttestationData) (int, int, error) {
	sourceEpoch, err := strconv.Atoi(att.Source.Epoch)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to convert source epoch to int")
		return 0, 0, err
	}
	targetEpoch, err := strconv.Atoi(att.Target.Epoch)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to convert target epoch to int")
		return 0, 0, err
	}
	return sourceEpoch, targetEpoch, err
}
