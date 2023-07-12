package artemis_trading_auxiliary

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (a *AuxiliaryTradingUtils) getNonce(ctx context.Context) (uint64, error) {
	nonce, err := a.w3a().GetNonce(ctx)
	if err != nil {
		log.Err(err).Msg("error getting nonce")
		return nonce, err
	}
	return nonce, nil
}
