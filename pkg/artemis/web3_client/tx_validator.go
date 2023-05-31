package web3_client

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

func (w *Web3Client) GetBlockTxs(ctx context.Context) (types.Transactions, error) {
	w.Dial()
	defer w.Close()
	block, err := w.C.BlockByNumber(ctx, nil)
	if err != nil {
		log.Err(err).Msg("failed to get nonce")
		return nil, err
	}
	return block.Transactions(), nil
}
