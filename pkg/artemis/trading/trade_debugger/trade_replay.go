package artemis_trade_debugger

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *TradeDebugger) getMevTx(ctx context.Context, txHash string, fromMempoolTx bool) (HistoricalAnalysisDebug, error) {
	if fromMempoolTx {
		return t.lookupMevMempoolTx(ctx, txHash)
	}
	return t.lookupMevTx(ctx, txHash)
}
func (t *TradeDebugger) Replay(ctx context.Context, txHash string, fromMempoolTx bool) error {
	mevTx, err := t.getMevTx(ctx, txHash, fromMempoolTx)
	if err != nil {
		return err
	}
	tf := mevTx.TradePrediction
	err = t.ResetAndSetupPreconditions(ctx, tf)
	if err != nil {
		return err
	}
	fmt.Println("ANALYZING tx: ", tf.Tx.Hash().String(), "at block: ", mevTx.GetBlockNumber())
	_, err = t.dat.GetSimUniswapClient().FrontRunTradeGetAmountsOut(&tf)
	if err != nil {
		err = t.analyzeDrift(ctx, tf.FrontRunTrade)
		return err
	}
	ac := t.dat.GetSimAuxClient()
	tf.FrontRunTrade.AmountOut = tf.FrontRunTrade.SimulatedAmountOut //  new(big.Int).SetInt64(0)
	ur, err := ac.GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.FrontRunTrade)
	if err != nil {
		return err
	}
	if ur == nil {
		return fmt.Errorf("ur is nil")
	}
	start := tf.FrontRunTrade.AmountIn
	num := 0
	denom := 1000
	for i := 0; i < 10; i++ {
		switch i {
		case 0:
			num = 0
		case 1:
			num = 1
		case 2:
			num = 10
		case 3:
			num = 50
		case 4:
			num = 100
		case 5:
			num = 200
		}
		tf.FrontRunTrade.AmountIn = artemis_eth_units.ApplyTransferTax(start, num, denom)
		err = t.dat.GetSimUniswapClient().InjectExecTradeV2SwapFromTokenToToken(ctx, ur, &tf.FrontRunTrade)
		if err == nil {
			log.Info().Interface("num", num).Msgf("Injected trade with amount in: %s", tf.FrontRunTrade.AmountIn.String())
			break
		}
	}
	if num == 0 {
		num = 1
		denom = 1
	}
	err = artemis_mev_models.UpdateERC20TokenTransferTaxInfo(ctx, artemis_autogen_bases.Erc20TokenInfo{
		Address:                tf.FrontRunTrade.AmountIn.String(),
		ProtocolNetworkID:      hestia_req_types.EthereumMainnetProtocolNetworkID,
		TransferTaxNumerator:   &num,
		TransferTaxDenominator: &denom,
	})
	if err != nil {
		return err
	}
	_, err = t.dat.GetSimUniswapClient().ExecTradeByMethod(&tf)
	if err != nil {
		return err
	}
	tf.SandwichTrade.AmountIn = tf.FrontRunTrade.SimulatedAmountOut
	_, err = t.dat.GetSimUniswapClient().SandwichTradeGetAmountsOut(&tf)
	if err != nil {
		err = t.analyzeDrift(ctx, tf.FrontRunTrade)
		return err
	}
	startBal, err := ac.CheckAuxERC20BalanceFromAddr(ctx, tf.SandwichTrade.AmountOutAddr.String())
	if err != nil {
		return err
	}
	tf.SandwichTrade.AmountOut = tf.SandwichTrade.SimulatedAmountOut
	ur, err = ac.GenerateTradeV2SwapFromTokenToToken(ctx, nil, &tf.SandwichTrade)
	if err != nil {
		return err
	}
	if ur == nil {
		return fmt.Errorf("ur is nil")
	}
	err = t.dat.GetSimUniswapClient().InjectExecTradeV2SwapFromTokenToToken(ctx, ur, &tf.SandwichTrade)
	if err != nil {
		return err
	}
	endBal, err := ac.CheckAuxERC20BalanceFromAddr(ctx, tf.SandwichTrade.AmountOutAddr.String())
	if err != nil {
		return err
	}
	fmt.Println("profit", artemis_eth_units.SubBigInt(endBal, startBal))
	fmt.Println("profitToken", tf.SandwichTrade.AmountOutAddr.String())
	//err = t.dat.GetSimUniswapClient().VerifyTradeResults(&tf)
	//if err != nil {
	//	return err
	//}
	return nil
}
