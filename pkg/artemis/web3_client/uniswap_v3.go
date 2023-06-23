package web3_client

import (
	"context"

	"github.com/rs/zerolog/log"
)

const (
	swapExactInputSingle    = "swapExactInputSingle"
	swapExactOutputSingle   = "swapExactOutputSingle"
	swapExactInputMultihop  = "swapExactInputMultihop"
	swapExactOutputMultihop = "swapExactOutputMultihop"
	multicall0              = "multicall0"
	multicall1              = "multicall1"
)

func (u *UniswapClient) ProcessUniswapV3RouterTxs(ctx context.Context, tx MevTx) {
	switch tx.MethodName {
	case swapExactInputSingle:
		inputs := &SwapExactInputSingleArgs{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode swap exact input single args")
			return
		}
		// convert, get pricing data, run bin search
	case swapExactOutputSingle:
		inputs := &SwapExactOutputSingleArgs{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode swap exact output single args")
			return
		}
		// convert, get pricing data, run bin search
	case multicall0, multicall1:
		inputs := &Multicall{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode multicall args")
			return
		}
	case swapExactInputMultihop:
	case swapExactOutputMultihop:
	}
}
