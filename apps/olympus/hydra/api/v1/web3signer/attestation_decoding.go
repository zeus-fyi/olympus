package hydra_eth2_web3signer

import (
	"context"

	"github.com/rs/zerolog/log"
	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
	ethereum_slashing_protection_watermarking "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/slashing_protection"
)

func CanSignAttestation(ctx context.Context, pubkey string, att consensys_eth2_openapi.AttestationData) error {
	err := ethereum_slashing_protection_watermarking.WatermarkAttestation(ctx, pubkey, att)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to check watermark for attestation")
		return err
	}
	return err
}
