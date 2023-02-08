package ethereum_slashing_protection_watermarking

import (
	"context"
	"github.com/rs/zerolog/log"
	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
	"strconv"
)

func WatermarkBlock(ctx context.Context, pubkey string, block consensys_eth2_openapi.BeaconBlock) error {
	slot, err := strconv.Atoi(block.Slot)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to convert source epoch to int")
		return err
	}

	// TODO
	tmpMockSlot := 0
	if slot > tmpMockSlot {

	}
	// TODO
	return err
}
