package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
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

/*
	tmTradingEnabled := TokenMap[to].TradingEnabled
	if tmTradingEnabled == nil {
		return errors.New("ActiveTrading: EntryTxFilter, erc20 at address not registered")
	}
	tradingEnabled := false
	tradingEnabled = *TokenMap[to].TradingEnabled
	if !tradingEnabled {
		return errors.New("ActiveTrading: EntryTxFilter, trading not enabled for this token")
	}
*/
