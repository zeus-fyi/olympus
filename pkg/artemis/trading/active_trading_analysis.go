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

/*
	MevSmartContractTxMapUniversalRouterNew: MevSmartContractTxMap{
		SmartContractAddr: UniswapUniversalRouterAddressNew,
		Abi:               artemis_oly_contract_abis.MustLoadNewUniversalRouterAbi(),
		Txs:               []MevTx{},
	},

	MevSmartContractTxMapUniversalRouterOld: MevSmartContractTxMap{
		SmartContractAddr: UniswapUniversalRouterAddressOld,
		Abi:               artemis_oly_contract_abis.MustLoadOldUniversalRouterAbi(),
		Txs:               []MevTx{},
	},

	MevSmartContractTxMapV2Router02: MevSmartContractTxMap{
		SmartContractAddr: UniswapV2Router02Address,
		Abi:               artemis_oly_contract_abis.MustLoadUniswapV2Router02ABI(),
		Txs:               []MevTx{},
		Filter:            &f,
	},

	MevSmartContractTxMapV2Router01: MevSmartContractTxMap{
		SmartContractAddr: UniswapV2Router01Address,
		Abi:               artemis_oly_contract_abis.MustLoadUniswapV2Router01ABI(),
		Txs:               []MevTx{},
		Filter:            &f,
	},

	MevSmartContractTxMapV3SwapRouterV1: MevSmartContractTxMap{
		SmartContractAddr: UniswapV3Router01Address,
		Abi:               artemis_oly_contract_abis.MustLoadUniswapV3Swap1RouterAbi(),
		Txs:               []MevTx{},
	},

	MevSmartContractTxMapV3SwapRouterV2: MevSmartContractTxMap{
		SmartContractAddr: UniswapV3Router02Address,
		Abi:               artemis_oly_contract_abis.MustLoadUniswapV3Swap2RouterAbi(),
		Txs:               []MevTx{},
	},
*/
func (a *ActiveTrading) ProcessTxs(ctx context.Context, mevTxs []web3_client.MevTx) ([]web3_client.TradeExecutionFlowJSON, error) {
	var tfSlice []web3_client.TradeExecutionFlowJSON
	for _, mevTx := range mevTxs {
		switch mevTx.Tx.To().String() {
		case artemis_trading_constants.UniswapUniversalRouterAddressOld:
			tf, err := a.RealTimeProcessUniversalRouterTx(ctx, mevTx)
			if err != nil {
				log.Err(err).Msg("error processing universal router tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		case artemis_trading_constants.UniswapUniversalRouterAddressNew:
			tf, err := a.RealTimeProcessUniversalRouterTx(ctx, mevTx)
			if err != nil {
				log.Err(err).Msg("error processing universal router tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		case artemis_trading_constants.UniswapV2Router01Address:
			tf, err := a.RealTimeProcessUniswapV2RouterTx(ctx, mevTx)
			if err != nil {
				log.Err(err).Msg("error processing v2_01 router tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		case artemis_trading_constants.UniswapV2Router02Address:
			tf, err := a.RealTimeProcessUniswapV2RouterTx(ctx, mevTx)
			if err != nil {
				log.Err(err).Msg("error processing v2_02 router tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		case artemis_trading_constants.UniswapV3Router01Address:
			tf, err := a.RealTimeProcessUniswapV3RouterTx(ctx, mevTx, UniswapV3Router01Abi, nil)
			if err != nil {
				log.Err(err).Msg("error processing v3_01 router tx")
				continue
			}
			tfSlice = append(tfSlice, tf...)
		case artemis_trading_constants.UniswapV3Router02Address:
			tf, err := a.RealTimeProcessUniswapV3RouterTx(ctx, mevTx, UniswapV3Router02Abi, nil)
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
