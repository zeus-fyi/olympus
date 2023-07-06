package artemis_trade_debugger

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
)

func (t *TradeDebugger) lookupMevTx(ctx context.Context, txHash string) ([]artemis_mev_models.HistoricalAnalysis, error) {
	mevTxs, merr := artemis_mev_models.SelectEthMevTxAnalysisByTxHash(ctx, txHash)
	if merr != nil {
		return nil, merr
	}
	return mevTxs, merr
}

func (t *TradeDebugger) getTxFromHash(ctx context.Context, txHash string) (*types.Transaction, error) {
	hash := common.HexToHash(txHash)
	tx, _, err := t.UniswapClient.Web3Client.GetTxByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (t *TradeDebugger) getRxFromHash(ctx context.Context, txHash string) (*types.Receipt, error) {
	hash := common.HexToHash(txHash)
	rx, err := t.UniswapClient.Web3Client.GetTxReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}
	return rx, nil
}
