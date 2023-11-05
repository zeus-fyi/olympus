package artemis_reporting

import (
	"fmt"

	"github.com/metachris/flashbotsrpc"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
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

func (s *ReportingTestSuite) TestInsertCallBundleResp() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	cr := flashbotsrpc.FlashbotsCallBundleResponse{
		BundleGasPrice:    "",
		BundleHash:        "0x",
		CoinbaseDiff:      "",
		EthSentToCoinbase: "",
		GasFees:           "",
		Results: []flashbotsrpc.FlashbotsCallBundleResult{
			{
				CoinbaseDiff:      "dd",
				EthSentToCoinbase: "s",
				FromAddress:       "d",
				GasFees:           "",
				GasPrice:          "",
				GasUsed:           1,
				ToAddress:         "",
				TxHash:            "0x",
				Value:             "",
				Error:             "",
				Revert:            "ss",
			},
		},
		StateBlockNumber: 1,
		TotalGasUsed:     1,
	}
	err := InsertCallBundleResp(ctx, "flashbots", 1, cr)
	s.Assert().Nil(err)
}
