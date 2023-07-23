package artemis_realtime_trading

import (
	"context"
	"errors"
	"fmt"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

var localCache = cache.New(cache.NoExpiration, cache.NoExpiration)

func CheckTokenRegistry(ctx context.Context, tokenAddress string, chainID int64) error {
	tmTradingEnabled := artemis_trading_cache.TokenMap[tokenAddress].TradingEnabled
	val, ok := localCache.Get(tokenAddress)
	if ok && val == true {
		return nil
	}
	if tmTradingEnabled == nil {
		tradeEnabled := false
		log.Info().Msgf("CheckTokenRegistry, erc20 at address %s not registered", tokenAddress)
		err := artemis_mev_models.InsertERC20TokenInfo(ctx, artemis_autogen_bases.Erc20TokenInfo{
			Address:           tokenAddress,
			ProtocolNetworkID: int(chainID),
			BalanceOfSlotNum:  -2, // -1 means balanceOf it wasn't cracked within 100 slots, -2 means cracking hasn't been attempted yet
			TradingEnabled:    &tradeEnabled,
		})
		if err != nil {
			log.Err(err).Msg("CheckTokenRegistry: InsertERC20TokenInfo")
			return errors.New("CheckTokenRegistry: erc20 at address %s not registered")
		}
		localCache.Set(tokenAddress, true, cache.NoExpiration)
	}
	return nil
}

func ApplyMaxTransferTax(ctx context.Context, tf *web3_client.TradeExecutionFlow) error {
	bn, berr := artemis_trading_cache.GetLatestBlock(context.Background())
	if berr != nil {
		log.Err(berr).Msg("failed to get latest block")
		return errors.New("ailed to get latest block")
	}
	tf.CurrentBlockNumber = artemis_eth_units.NewBigIntFromUint(bn)
	tokenOne := tf.UserTrade.AmountInAddr.String()
	tokenTwo := tf.UserTrade.AmountOutAddr.String()
	if tokenOne == artemis_trading_constants.ZeroAddress && tokenTwo == artemis_trading_constants.ZeroAddress {
		log.Warn().Str("txHash", tf.Tx.Hash().String()).Str("tradeMethod", tf.Trade.TradeMethod).Interface("tf.UserTrade", tf.UserTrade).Str("toAddr", tf.Tx.To().String()).Msg("dat: ApplyMaxTransferTax, tokenOne and tokenTwo are zero address")
		return errors.New("dat: ApplyMaxTransferTax, tokenOne and tokenTwo are zero address")
	}
	go func(ctx context.Context, tokenA, tokenB string) {
		if tokenA != artemis_trading_constants.ZeroAddress && tokenA != artemis_trading_constants.WETH9ContractAddress {
			err := CheckTokenRegistry(ctx, tokenA, hestia_req_types.EthereumMainnetProtocolNetworkID)
			if err != nil {
				log.Err(err).Msg("CheckTokenRegistry: failed to check token registry")
			}
		}
		if tokenB != artemis_trading_constants.ZeroAddress && tokenB != artemis_trading_constants.WETH9ContractAddress {
			err := CheckTokenRegistry(ctx, tokenB, hestia_req_types.EthereumMainnetProtocolNetworkID)
			if err != nil {
				log.Err(err).Msg("CheckTokenRegistry: failed to check token registry")
			}
		}
	}(context.Background(), tokenOne, tokenTwo)

	maxNum, maxDen := 0, 1
	if info, ok := artemis_trading_cache.TokenMap[tokenOne]; ok {
		if info.TransferTaxNumerator == nil || info.TransferTaxDenominator == nil {
			fmt.Println("token not found in cache")
		} else {
			den := info.TransferTaxDenominator
			num := info.TransferTaxNumerator
			if den != nil && num != nil {
				fmt.Println("token: ", tokenOne, "transferTax: num: ", *num, "den: ", *den)
				if *num > maxNum {
					maxNum = *num
					maxDen = *den
				}
			}
		}
	}
	if info, ok := artemis_trading_cache.TokenMap[tokenTwo]; ok {
		if info.TransferTaxNumerator == nil || info.TransferTaxDenominator == nil {
			fmt.Println("token not found in cache")
		} else {
			den := info.TransferTaxDenominator
			num := info.TransferTaxNumerator
			if den != nil && num != nil {
				fmt.Println("token: ", tokenTwo, "tradingTax: num: ", *num, "den: ", *den)
				if *num > maxNum {
					maxNum = *num
					maxDen = *den
				}
			}
		}
	}

	if !tf.AreAllTradesValid() {
		log.Warn().Msg("ApplyMaxTransferTax: trades are not valid")
		return errors.New("ApplyMaxTransferTax: trades are not valid")
	}
	fmt.Println("maxNum: ", maxNum, "maxDen: ", maxDen)

	if maxNum == 0 {
		amountOutStartFrontRun := tf.FrontRunTrade.AmountOut
		amountOutStartSandwich := tf.SandwichTrade.AmountOut

		adjAmountOutFrontRun := artemis_eth_units.ApplyTransferTax(amountOutStartFrontRun, 5, 1000)
		tf.FrontRunTrade.AmountOut = adjAmountOutFrontRun

		tf.SandwichTrade.AmountIn = tf.FrontRunTrade.AmountOut

		adjAmountOutSandwich := artemis_eth_units.ApplyTransferTax(amountOutStartSandwich, 10, 1000)
		tf.SandwichTrade.AmountOut = adjAmountOutSandwich
		tf.SandwichPrediction.ExpectedProfit = adjAmountOutSandwich
		if !tf.AreAllTradesValid() {
			log.Warn().Msg("ApplyMaxTransferTax: trades are not valid")
			return errors.New("ApplyMaxTransferTax: trades are not valid")
		}
		log.Info().Str("txHash", tf.Tx.Hash().String()).Uint64("bn", tf.CurrentBlockNumber.Uint64()).Str("profitTokenAddress", tf.SandwichTrade.AmountOutAddr.String()).Interface("sellAmount", tf.SandwichPrediction.SellAmount).Interface("tf.SandwichPrediction.ExpectedProfit", tf.SandwichPrediction.ExpectedProfit).Str("tf.SandwichTrade.AmountOut", tf.SandwichTrade.AmountOut.String()).Msg("ApplyMaxTransferTax: acceptable after tax")
		return nil
	}
	amountOutStartFrontRun := tf.FrontRunTrade.AmountOut
	amountOutStartSandwich := tf.SandwichTrade.AmountOut

	adjAmountOutFrontRun := artemis_eth_units.ApplyTransferTax(amountOutStartFrontRun, maxNum, maxDen)
	tf.FrontRunTrade.AmountOut = adjAmountOutFrontRun

	tf.SandwichTrade.AmountIn = tf.FrontRunTrade.AmountOut
	if maxDen == 1000 {
		maxNum += 5
	}
	adjAmountOutSandwich := artemis_eth_units.ApplyTransferTax(amountOutStartSandwich, maxNum, maxDen)
	tf.SandwichTrade.AmountOut = adjAmountOutSandwich
	tf.SandwichPrediction.ExpectedProfit = adjAmountOutSandwich

	if !tf.AreAllTradesValid() {
		log.Info().Str("txHash", tf.Tx.Hash().String()).Uint64("bn", tf.CurrentBlockNumber.Uint64()).Str("profitTokenAddress", tf.SandwichTrade.AmountOutAddr.String()).Str("startingSandwichOut", amountOutStartSandwich.String()).Interface("sellAmount", tf.SandwichPrediction.SellAmount).Interface("tf.SandwichPrediction.ExpectedProfit", tf.SandwichPrediction.ExpectedProfit).Str("tf.SandwichTrade.AmountOut", tf.SandwichTrade.AmountOut.String()).Msg("ApplyMaxTransferTax: trade not acceptable after tax")
		return errors.New("ApplyMaxTransferTax: trades are not valid")
	}
	log.Info().Str("txHash", tf.Tx.Hash().String()).Uint64("bn", tf.CurrentBlockNumber.Uint64()).Str("profitTokenAddress", tf.SandwichTrade.AmountOutAddr.String()).Str("startingSandwichOut", amountOutStartSandwich.String()).Interface("sellAmount", tf.SandwichPrediction.SellAmount).Interface("tf.SandwichPrediction.ExpectedProfit", tf.SandwichPrediction.ExpectedProfit).Str("tf.SandwichTrade.AmountOut", tf.SandwichTrade.AmountOut.String()).Msg("ApplyMaxTransferTax: acceptable after tax")
	return nil
}
