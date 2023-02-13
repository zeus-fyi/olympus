package hydra_eth2_web3signer

import (
	"context"
	"github.com/rs/zerolog/log"
	"strconv"

	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
	ethereum_slashing_protection_watermarking "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/slashing_protection"
)

func CanSignAttestation(ctx context.Context, pubkey string, att consensys_eth2_openapi.AttestationData) error {
	sourceEpoch, err := strconv.Atoi(att.Source.Epoch)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to convert source epoch to int")
		return err
	}
	targetEpoch, err := strconv.Atoi(att.Target.Epoch)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to convert target epoch to int")
		return err
	}
	err = ethereum_slashing_protection_watermarking.WatermarkAttestation(ctx, pubkey, sourceEpoch, targetEpoch)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to check watermark for attestation")
		return err
	}
	return err
}
