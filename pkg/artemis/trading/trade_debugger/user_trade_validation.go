package artemis_trade_debugger

import (
	"context"
	"fmt"

	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (t *TradeDebugger) analyzeUserTrade(ctx context.Context, tf *web3_client.TradeExecutionFlow) error {
	if tf.Tx == nil {
		return fmt.Errorf("tx is nil")
	}
	from, err := web3_actions.GetSender(tf.Tx)
	if err != nil {
		return err
	}
	nonceAt, err := t.UniswapClient.Web3Client.NonceAt(ctx, from, tf.CurrentBlockNumber)
	if err != nil {
		return err
	}

	pendingNonce, err := t.UniswapClient.Web3Client.PendingNonce(ctx, from)
	if err != nil {
		return err
	}
	fmt.Println("from ", from)
	fmt.Println("nonceAt ", nonceAt)
	fmt.Println("exp nonce ", pendingNonce)
	fmt.Println("tx nonce ", tf.Tx.Nonce())
	fmt.Println("value ", tf.Tx.Value())
	fmt.Println("gasPrice ", tf.Tx.GasPrice())
	fmt.Println("gas ", tf.Tx.Gas())
	return nil
}
