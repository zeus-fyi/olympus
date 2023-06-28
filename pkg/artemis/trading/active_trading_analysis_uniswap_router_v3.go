package artemis_realtime_trading

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

const (
	multicall               = "multicall"
	swapExactInputSingle    = "swapExactInputSingle"
	swapExactOutputSingle   = "swapExactOutputSingle"
	swapExactInputMultihop  = "swapExactInputMultihop"
	swapExactOutputMultihop = "swapExactOutputMultihop"
	exactInput              = "exactInput"
	exactOutput             = "exactOutput"
)

func (a *ActiveTrading) RealTimeProcessUniswapV3RouterTx(ctx context.Context, tx web3_client.MevTx, abiFile *abi.ABI, filter *strings_filter.FilterOpts) {
	if tx.Tx.To() == nil {
		return
	}
	toAddr := tx.Tx.To().String()
	if strings.HasPrefix(tx.MethodName, multicall) {
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, multicall)
		inputs := &web3_client.Multicall{}
		err := inputs.Decode(ctx, tx.Args)
		if err != nil {
			log.Err(err).Msg("failed to decode multicall args")
			return
		}
		for _, data := range inputs.Data {
			mn, args, derr := web3_client.DecodeTxData(ctx, data, abiFile, filter)
			if derr != nil {
				log.Err(derr).Msg("failed to decode tx data")
				continue
			}
			newTx := tx
			newTx.MethodName = mn
			newTx.Args = args
			a.processUniswapV3Txs(ctx, newTx)
		}
	} else {
		a.processUniswapV3Txs(ctx, tx)
	}
	return
}

func (a *ActiveTrading) processUniswapV3Txs(ctx context.Context, tx web3_client.MevTx) {
	if tx.Tx.To() == nil {
		return
	}
	toAddr := tx.Tx.To().String()
	switch tx.MethodName {
	case exactInput:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, exactInput)
	case exactOutput:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, exactOutput)
	case swapExactInputSingle:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactInputSingle)
	case swapExactOutputSingle:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactOutputSingle)
	case swapExactTokensForTokens:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactTokensForTokens)
	case swapExactInputMultihop:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactInputMultihop)
	case swapExactOutputMultihop:
		a.m.TxFetcherMetrics.TransactionGroup(toAddr, swapExactOutputMultihop)
	}
}
