package artemis_realtime_trading

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

/*
  adding in other filters here
	  - filter by token
	  - filter by profit
	  - filter by risk score
	  - adds sourcing of new blocks
*/

func (a *ActiveTrading) EntryTxFilter(ctx context.Context, tx *types.Transaction) error {
	if tx.To() == nil {
		return errors.New("ActiveTrading: EntryTxFilter, tx.To() is nil")
	}
	if tx.ChainId() == nil {
		return errors.New("ActiveTrading: EntryTxFilter, tx.ChainId() is nil")
	}
	return nil
}

func (a *ActiveTrading) SimTxFilter(ctx context.Context, tfSlice []web3_client.TradeExecutionFlowJSON) error {
	var addresses []accounts.Address
	for _, tf := range tfSlice {
		addresses = append(addresses, tf.UserTrade.AmountInAddr)
		addresses = append(addresses, tf.UserTrade.AmountOutAddr)
	}
	for _, addr := range addresses {
		if artemis_trading_cache.TokenMap[addr.String()].BalanceOfSlotNum < 0 {
			return errors.New("ActiveTrading: EntryTxFilter, balanceOf not cracked yet")
		}
		num := artemis_trading_cache.TokenMap[addr.String()].TransferTaxNumerator
		den := artemis_trading_cache.TokenMap[addr.String()].TransferTaxDenominator
		if num == nil || den == nil {
			return errors.New("ActiveTrading: EntryTxFilter, transfer tax not set")
		}
		if *num == 0 || *den == 0 {
			return errors.New("ActiveTrading: EntryTxFilter, transfer tax not set")
		}
	}
	return nil
}

func (a *ActiveTrading) ActiveTradingFilter(ctx context.Context, tf web3_client.TradeExecutionFlowJSON) error {
	if tf.UserTrade.AmountInAddr.String() != artemis_trading_constants.WETH9ContractAddressAccount.String() {
		return errors.New("ActiveTrading: ActiveTradingFilter: only WETH is supported as amountIn for now")
	}
	err := a.TradingEnabledFilter(ctx, tf.UserTrade.AmountOutAddr)
	if err != nil {
		return err
	}
	switch tf.Trade.TradeMethod {
	case artemis_trading_constants.SwapExactTokensForTokens:
	case artemis_trading_constants.V2SwapExactIn, artemis_trading_constants.V2SwapExactOut:
	default:
		return fmt.Errorf("ActiveTrading: ActiveTradingFilter: %s method not supported for now", tf.Trade.TradeMethod)
	}
	return nil
}

func (a *ActiveTrading) TradingEnabledFilter(ctx context.Context, address accounts.Address) error {
	return errors.New("ActiveTrading: ActiveTradingFilter: trading not enabled for token")
}
