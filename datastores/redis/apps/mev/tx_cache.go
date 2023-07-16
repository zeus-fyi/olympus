package redis_mev

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

func (m *MevCache) AddTxHashCache(ctx context.Context, txHash string, ttl time.Duration) error {
	//txHash := tx.Hash().String()
	statusCmd := m.Set(ctx, txHash, true, ttl)
	if statusCmd.Err() != nil {
		log.Ctx(ctx).Err(statusCmd.Err()).Msgf("SetTxHashCache: %s", txHash)
		return statusCmd.Err()
	}
	return nil
}

func (m *MevCache) DoesTxExist(ctx context.Context, txHash string) (bool, error) {
	log.Info().Msgf("DoesTxExist: %s", txHash)
	err := m.Get(ctx, txHash).Err()
	switch {
	case err == redis.Nil:
		log.Info().Msgf("DoesTxExist: tx hash not previously seen")
		return false, nil
	case err != nil:
		log.Err(err).Msgf("DoesTxExist: %s", txHash)
		return false, err
	}
	return true, err
}

func (m *MevCache) DeleteTx(ctx context.Context, txHash string) error {
	err := m.Del(ctx, txHash)
	if err != nil {
		log.Err(err.Err()).Msgf("DeleteTx: %s", txHash)
		return err.Err()
	}
	return err.Err()
}
