package artemis_trade_debugger

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type TradeDebugger struct {
	UniswapClient *web3_client.UniswapClient
	ActiveTrading artemis_realtime_trading.ActiveTrading
}

func NewTradeDebugger(a artemis_realtime_trading.ActiveTrading, u *web3_client.UniswapClient) TradeDebugger {
	return TradeDebugger{
		ActiveTrading: a,
		UniswapClient: u,
	}
}

func (t *TradeDebugger) GetTxFromHash(ctx context.Context, txHash string) error {
	hash := common.HexToHash(txHash)
	rx, _, err := t.UniswapClient.Web3Client.GetTxByHash(ctx, hash)
	if err != nil {
		return err
	}
	fmt.Println(rx)
	return nil

}
