package artemis_reporting

import (
	"fmt"
	"strings"
	"unicode"

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
		BundleGasPrice: "",
		BundleHash:     "0x",

		StateBlockNumber: 1,
		TotalGasUsed:     1,
	}
	err := InsertCallBundleResp(ctx, "flashbots", 1, cr)
	s.Assert().Nil(err)
}

func (s *ReportingTestSuite) TestTTT() {
	str := `\u0008�y�\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000 \u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0014TRANSFER_FROM_FAILED\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000"}],"`
	cleanStr := removeInvalidUnicode(str)

	s.Assert().NotEqual(str, cleanStr)
	fmt.Println("cleanStr", cleanStr)
}

func removeInvalidUnicode(input string) string {
	var sb strings.Builder
	for _, r := range input {
		if r == unicode.ReplacementChar || !unicode.IsPrint(r) {
			// Skip this rune
			continue
		}
		sb.WriteRune(r)
	}
	return sb.String()
}
