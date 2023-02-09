package ethereum_slashing_protection_watermarking

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
)

func WatermarkBlock(ctx context.Context, pubkey string, proposedSlot int) error {
	prevSlot, err := FetchLastSignedBlockSlot(ctx, pubkey)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to fetch last signed block slot")
		return err
	}
	if proposedSlot <= prevSlot {
		log.Ctx(ctx).Warn().Msgf("proposed slot %d less than or equal to a previous block slot %d", proposedSlot, prevSlot)
		return errors.New("proposed slot less than or equal to a previous block slot")
	}
	return nil
}

func FetchLastSignedBlockSlot(ctx context.Context, pubkey string) (prevSlot int, err error) {
	// TODO
	return 0, nil
}
