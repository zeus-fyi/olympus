package artemis_realtime_trading

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) ProcessTxs(ctx context.Context, mevTxs []web3_client.MevTx, m *metrics_trading.TradingMetrics, w3a web3_actions.Web3Actions) ([]web3_client.TradeExecutionFlowJSON, error) {
	var tfSlice []web3_client.TradeExecutionFlowJSON
	for _, mevTx := range mevTxs {
		switch mevTx.Tx.To().String() {
		case artemis_trading_constants.UniswapUniversalRouterAddressOld:
			tf, err := RealTimeProcessUniversalRouterTx(ctx, mevTx, m, w3a)
			if err != nil {
				log.Err(err).Msg("error processing universal router tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		case artemis_trading_constants.UniswapUniversalRouterAddressNew:
			tf, err := RealTimeProcessUniversalRouterTx(ctx, mevTx, m, w3a)
			if err != nil {
				log.Err(err).Msg("error processing universal router tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		case artemis_trading_constants.UniswapV2Router01Address:
			tf, err := RealTimeProcessUniswapV2RouterTx(ctx, mevTx, m, w3a)
			if err != nil {
				log.Err(err).Msg("error processing v2_01 router tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		case artemis_trading_constants.UniswapV2Router02Address:
			tf, err := RealTimeProcessUniswapV2RouterTx(ctx, mevTx, m, w3a)
			if err != nil {
				log.Err(err).Msg("error processing v2_02 router tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		case artemis_trading_constants.UniswapV3Router01Address:
			tf, err := RealTimeProcessUniswapV3RouterTx(ctx, mevTx, UniswapV3Router01Abi, nil, m, w3a)
			if err != nil {
				log.Err(err).Msg("error processing v3_01 router tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		case artemis_trading_constants.UniswapV3Router02Address:
			tf, err := RealTimeProcessUniswapV3RouterTx(ctx, mevTx, UniswapV3Router02Abi, nil, m, w3a)
			if err != nil {
				log.Err(err).Msg("error processing v3_02 router tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		}
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

func ApplyMaxTransferTax(tf *web3_client.TradeExecutionFlowJSON) {
	tokenOne := tf.UserTrade.AmountInAddr.String()
	tokenTwo := tf.UserTrade.AmountOutAddr.String()
	maxNum, maxDen := 1, 1
	if info, ok := artemis_trading_cache.TokenMap[tokenOne]; ok {
		den := info.TransferTaxDenominator
		num := info.TransferTaxNumerator
		if den != nil && num != nil {
			fmt.Println("token: ", tokenOne, "transferTax: num: ", *num, "den: ", *den)

			if *num > maxNum {
				maxNum = *num
				maxDen = *den
			}
		} else {
			fmt.Println("token not found in cache")
		}
	}
	if info, ok := artemis_trading_cache.TokenMap[tokenTwo]; ok {
		den := info.TransferTaxDenominator
		num := info.TransferTaxNumerator
		if den != nil && num != nil {
			fmt.Println("token: ", tokenTwo, "tradingTax: num: ", *num, "den: ", *den)
			if *num > maxNum {
				maxNum = *num
				maxDen = *den
			}
		} else {
			fmt.Println("token not found in cache")
		}
	}
	amountOutStartFrontRun := artemis_eth_units.NewBigIntFromStr(tf.FrontRunTrade.AmountOut)
	amountOutStartSandwich := artemis_eth_units.NewBigIntFromStr(tf.SandwichTrade.AmountOut)

	adjAmountOutFrontRun := artemis_eth_units.ApplyTransferTax(amountOutStartFrontRun, maxNum, maxDen)
	tf.FrontRunTrade.AmountOut = adjAmountOutFrontRun.String()

	tf.SandwichTrade.AmountIn = tf.FrontRunTrade.AmountOut
	adjAmountOutSandwich := artemis_eth_units.ApplyTransferTax(amountOutStartSandwich, maxNum+30, maxDen)
	tf.SandwichTrade.AmountOut = adjAmountOutSandwich.String()
	tf.SandwichPrediction.ExpectedProfit = adjAmountOutSandwich.String()
	fmt.Println("maxNum: ", maxNum, "maxDen: ", maxDen)
}
