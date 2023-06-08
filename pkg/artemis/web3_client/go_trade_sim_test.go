package web3_client

import (
	"encoding/json"
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

/*
	modEven :=  tokenSellAmount.Mod(tokenSellAmount, big.NewInt(2))
	if modEven.String() == "0" {
		tokenSellAmount = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
	} else {
		tokenSellAmount = tokenSellAmount.Add(tokenSellAmount, big.NewInt(1))
		tokenSellAmount = tokenSellAmount.Div(tokenSellAmount, big.NewInt(2))
	}
*/

func (s *Web3ClientTestSuite) TestJson() {
	//amountOutMin, _ := new(big.Int).SetString("746627207819418433569734379647", 10)
	//te := TradeExecutionFlowJSON{
	//	CurrentBlockNumber: nil,
	//	Tx:                 nil,
	//	Trade: Trade{
	//		TradeMethod:                    "swapExactETHForTokens",
	//		SwapETHForExactTokensParams:    nil,
	//		SwapTokensForExactTokensParams: nil,
	//		SwapExactTokensForTokensParams: nil,
	//		SwapExactETHForTokensParams: &SwapExactETHForTokensParams{
	//			AmountOutMin: amountOutMin,
	//		},
	//		SwapExactTokensForETHParams: nil,
	//		SwapTokensForExactETHParams: nil,
	//	},
	//	InitialPair:        UniswapV2Pair{},
	//	FrontRunTrade:      TradeOutcome{},
	//	SandwichTrade:      TradeOutcome{},
	//	SandwichPrediction: SandwichTradePrediction{},
	//}
	//
	//b, _ := json.Marshal(te)
	//fmt.Println(string(b))
}

func (s *Web3ClientTestSuite) TestTradeSim() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	uni := InitUniswapClient(ctx, s.MainnetWeb3User)
	mevTxs, err := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17326677)
	s.Require().Nil(err)
	s.Require().NotEmpty(mevTxs)

	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlowJSON{}
		b := []byte(mevTx.TxFlowPrediction)
		berr := json.Unmarshal(b, &tf)
		s.Require().Nil(berr)

		if tf.FrontRunTrade.AmountIn == "" {
			continue
		}
		fmt.Println(tf.Trade.TradeMethod)
		b, berr = json.MarshalIndent(tf, "", "  ")
		s.Require().Nil(berr)
		fmt.Println(string(b))
		tfConv := tf.ConvertToBigIntType()
		executedProfit := uni.TradeSimStep(tfConv)
		fmt.Println("binary search sell amount", tf.SandwichPrediction.SellAmount)
		fmt.Println("binary search max profit", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("profit from execution path", executedProfit.String())
		fmt.Println("profit token type", tf.SandwichTrade.AmountOutAddr.String())
		s.Assert().Equal(tf.SandwichPrediction.ExpectedProfit, executedProfit.String())

		//sellAmount, maxProfit := uni.TradeSim(tfConv)
		//fmt.Println("linear search sell amount", sellAmount.String())
		//fmt.Println("linear search max profit", maxProfit.String())
		//tf.FrontRunTrade.AmountIn = sellAmount.String()
	}
}
