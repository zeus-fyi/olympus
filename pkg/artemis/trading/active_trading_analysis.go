package artemis_realtime_trading

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *ActiveTrading) ProcessTxs(ctx context.Context) ([]web3_client.TradeExecutionFlowJSON, error) {
	var tfSlice []web3_client.TradeExecutionFlowJSON
	for _, mevTx := range a.a.U.MevSmartContractTxMapUniversalRouterOld.Txs {
		tf, err := a.RealTimeProcessUniversalRouterTx(ctx, mevTx)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.a.U.MevSmartContractTxMapUniversalRouterNew.Txs {
		tf, err := a.RealTimeProcessUniversalRouterTx(ctx, mevTx)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.a.U.MevSmartContractTxMapV2Router01.Txs {
		tf, err := a.RealTimeProcessUniswapV2RouterTx(ctx, mevTx)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.a.U.MevSmartContractTxMapV2Router02.Txs {
		tf, err := a.RealTimeProcessUniswapV2RouterTx(ctx, mevTx)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.a.U.MevSmartContractTxMapV3SwapRouterV2.Txs {
		tf, err := a.RealTimeProcessUniswapV3RouterTx(ctx, mevTx, a.a.U.MevSmartContractTxMapV3SwapRouterV2.Abi, a.a.U.MevSmartContractTxMapV3SwapRouterV2.Filter)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}
	for _, mevTx := range a.a.U.MevSmartContractTxMapV3SwapRouterV1.Txs {
		tf, err := a.RealTimeProcessUniswapV3RouterTx(ctx, mevTx, a.a.U.MevSmartContractTxMapV3SwapRouterV1.Abi, a.a.U.MevSmartContractTxMapV3SwapRouterV1.Filter)
		if err != nil {
			return nil, err
		}
		tfSlice = append(tfSlice, tf...)
	}

	for _, tf := range tfSlice {
		baseTx, err := tf.Tx.ConvertToTx()
		if err != nil {
			log.Err(err).Msg("ActiveTrading: EntryTxFilter, ConvertToTx")
			return nil, err
		}
		chainID := baseTx.ChainId().Int64()
		err = CheckTokenRegistry(ctx, tf.UserTrade.AmountInAddr.String(), chainID)
		if err != nil {
			log.Err(err).Msg("ActiveTrading: EntryTxFilter, CheckTokenRegistry")
			return nil, err
		}
		err = CheckTokenRegistry(ctx, tf.UserTrade.AmountOutAddr.String(), chainID)
		if err != nil {
			log.Err(err).Msg("ActiveTrading: EntryTxFilter, CheckTokenRegistry")
			return nil, err
		}
	}
	return tfSlice, nil
}

func CheckTokenRegistry(ctx context.Context, tokenAddress string, chainID int64) error {
	tmTradingEnabled := artemis_trading_cache.TokenMap[tokenAddress].TradingEnabled
	if tmTradingEnabled == nil {
		tradeEnabled := false
		log.Info().Msgf("ActiveTrading: EntryTxFilter, erc20 at address %s not registered", tokenAddress)
		err := artemis_mev_models.InsertERC20TokenInfo(ctx, artemis_autogen_bases.Erc20TokenInfo{
			Address:           tokenAddress,
			ProtocolNetworkID: int(chainID),
			BalanceOfSlotNum:  -2, // -1 means balanceOf it wasn't cracked within 100 slots, -2 means cracking hasn't been attempted yet
			TradingEnabled:    &tradeEnabled,
		})
		if err != nil {
			log.Err(err).Msg("ActiveTrading: EntryTxFilter, InsertERC20TokenInfo")
			return errors.New("ActiveTrading: EntryTxFilter, erc20 at address %s not registered")
		}
	}
	return nil
}
