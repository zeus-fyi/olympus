package artemis_realtime_trading

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) ProcessTxs(ctx context.Context) ([]web3_client.TradeExecutionFlowJSON, error) {
	var tfSlice []web3_client.TradeExecutionFlowJSON
	for _, mevTx := range a.GetUniswapClient().MevSmartContractTxMapUniversalRouterOld.Txs {
		tf, err := a.RealTimeProcessUniversalRouterTx(ctx, mevTx)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.GetUniswapClient().MevSmartContractTxMapUniversalRouterNew.Txs {
		tf, err := a.RealTimeProcessUniversalRouterTx(ctx, mevTx)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.GetUniswapClient().MevSmartContractTxMapV2Router01.Txs {
		tf, err := a.RealTimeProcessUniswapV2RouterTx(ctx, mevTx)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.GetUniswapClient().MevSmartContractTxMapV2Router02.Txs {
		tf, err := a.RealTimeProcessUniswapV2RouterTx(ctx, mevTx)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.GetUniswapClient().MevSmartContractTxMapV3SwapRouterV2.Txs {
		tf, err := a.RealTimeProcessUniswapV3RouterTx(ctx, mevTx, a.GetUniswapClient().MevSmartContractTxMapV3SwapRouterV2.Abi, a.GetUniswapClient().MevSmartContractTxMapV3SwapRouterV2.Filter)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.GetUniswapClient().MevSmartContractTxMapV3SwapRouterV1.Txs {
		tf, err := a.RealTimeProcessUniswapV3RouterTx(ctx, mevTx, a.GetUniswapClient().MevSmartContractTxMapV3SwapRouterV1.Abi, a.GetUniswapClient().MevSmartContractTxMapV3SwapRouterV1.Filter)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}

	var postFilter []web3_client.TradeExecutionFlowJSON
	for _, tf := range tfSlice {
		key := fmt.Sprintf("%s-tf", tf.Tx.Hash)
		_, ok := txCache.Get(key)
		if ok {
			log.Info().Msgf("dat: EntryTxFilter, tx already in cache, hash: %s", tf.Tx.Hash)
			continue
		}
		if tf.SandwichPrediction.ExpectedProfit == "" {
			continue
		}
		if tf.SandwichPrediction.ExpectedProfit == "0" {
			continue
		}
		baseTx, err := tf.Tx.ConvertToTx()
		if err != nil {
			log.Err(err).Msg("dat: EntryTxFilter, ConvertToTx")
			return nil, err
		}
		chainID := baseTx.ChainId().Int64()
		if tf.UserTrade.AmountInAddr.String() == artemis_trading_constants.WETH9ContractAddressAccount.String() {
			if tf.SandwichPrediction.ExpectedProfit != "0" {
				log.Info().Msgf("dat: EntryTxFilter, WETH9ContractAddressAccount, expected profit: %s, amountOutAddr %s", tf.SandwichPrediction.ExpectedProfit, tf.FrontRunTrade.AmountOutAddr.String())
			}
		}
		log.Info().Interface("userTrade", tf.UserTrade)
		log.Info().Interface("sandwichPrediction", tf.SandwichPrediction)
		err = CheckTokenRegistry(ctx, tf.UserTrade.AmountInAddr.String(), chainID)
		if err != nil {
			log.Err(err).Msg("dat: EntryTxFilter, CheckTokenRegistry")
			return nil, err
		}
		err = CheckTokenRegistry(ctx, tf.UserTrade.AmountOutAddr.String(), chainID)
		if err != nil {
			log.Err(err).Msg("dat: EntryTxFilter, CheckTokenRegistry")
			return nil, err
		}
		txCache.Set(key, tf, time.Hour*24)
		postFilter = append(postFilter, tf)
	}
	return postFilter, nil
}

func CheckTokenRegistry(ctx context.Context, tokenAddress string, chainID int64) error {
	tmTradingEnabled := artemis_trading_cache.TokenMap[tokenAddress].TradingEnabled
	if tmTradingEnabled == nil {
		tradeEnabled := false
		log.Info().Msgf("dat: EntryTxFilter, erc20 at address %s not registered", tokenAddress)
		err := artemis_mev_models.InsertERC20TokenInfo(ctx, artemis_autogen_bases.Erc20TokenInfo{
			Address:           tokenAddress,
			ProtocolNetworkID: int(chainID),
			BalanceOfSlotNum:  -2, // -1 means balanceOf it wasn't cracked within 100 slots, -2 means cracking hasn't been attempted yet
			TradingEnabled:    &tradeEnabled,
		})
		if err != nil {
			log.Err(err).Msg("dat: EntryTxFilter, InsertERC20TokenInfo")
			return errors.New("dat: EntryTxFilter, erc20 at address %s not registered")
		}
	}
	return nil
}
