package artemis_reporting

import (
	"fmt"
)

func (s *ReportingTestSuite) TestGetBundlesProfitHistory() {
	bg, err := GetBundlesProfitHistory(ctx, 0, 1)
	s.Assert().Nil(err)
	s.Assert().NotNil(bg)
	for bundleHash, b := range bg.Map {
		for _, bundleTx := range b {
			s.Require().Equal(bundleHash, bundleTx.BundleHash)
			s.Require().NotNil(bundleTx.TradeExecutionFlow)
			fmt.Println("bundleTx.EthTx.TxHash", bundleTx.EthTx.TxHash, "EthTxReceipts.EffectiveGasPrice", bundleTx.EffectiveGasPrice,
				"bundleTx.EthTxGas.GasTipCap", bundleTx.EthTxGas.GasTipCap, "bundleTx.EthTxGas.GasFeeCap", bundleTx.EthTxGas.GasFeeCap,
				"bundleTx.EthTxGas.GasLimit", bundleTx.EthTxGas.GasLimit)

			fmt.Println("bundleTx.EthMevBundleProfit.Profit", bundleTx.EthMevBundleProfit.Profit,
				"bundleTx.EthMevBundleProfit.Costs", bundleTx.EthMevBundleProfit.Costs,
				"bundleTx.EthMevBundleProfit.RevenuePrediction", bundleTx.EthMevBundleProfit.RevenuePrediction,
				"bundleTx.EthMevBundleProfit.RevenuePredictionSkew", bundleTx.EthMevBundleProfit.RevenuePredictionSkew)
		}
	}
}
