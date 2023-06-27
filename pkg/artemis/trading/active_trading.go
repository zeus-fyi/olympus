package artemis_realtime_trading

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

/*
	part 0. trade filter

  - should filter trades on multiple criteria
	- eg. profit denomination etc
	- eg. token risk score

	part 1. fast parallel processing of txs

  - should decode tx
  - should set up balances
  - should minimize rpc calls

	part 2. bundle txs into a single tx

  - should bundle txs into a single tx from processed ones

Overall:
	- should capture sim runtime metrics
	- should send trade metrics to prometheus
*/

type ActiveTrading struct {
	u *web3_client.UniswapClient
}

func NewActiveTradingModule(u *web3_client.UniswapClient) ActiveTrading {
	return ActiveTrading{u}
}

func (a *ActiveTrading) IngestTx(ctx context.Context, tx *types.Transaction) {
	a.ProcessTx(ctx, tx)
}
