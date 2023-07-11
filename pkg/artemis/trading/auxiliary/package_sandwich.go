package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) PackageSandwich(ctx context.Context, tf *web3_client.TradeExecutionFlow) (*web3_client.UniversalRouterExecCmd, error) {
	if tf == nil || tf.Tx == nil {
		return nil, errors.New("tf is nil")
	}
	// front run
	ur, err := a.GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.FrontRunTrade)
	if err != nil {
		return nil, err
	}
	// user trade
	err = a.AddTxToBundleGroup(ctx, tf.Tx)
	if err != nil {
		return nil, err
	}
	// sandwich trade
	ur, err = a.GenerateTradeV2SwapFromTokenToToken(ctx, ur, &tf.SandwichTrade)
	if err != nil {
		return nil, err
	}
	return ur, nil
}
