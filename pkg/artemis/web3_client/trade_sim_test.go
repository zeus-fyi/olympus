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

func (s *Web3ClientTestSuite) TestTradeSim() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	uni := InitUniswapV2Client(ctx, s.MainnetWeb3User)
	mevTxs, err := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17276013)
	s.Require().Nil(err)
	s.Require().NotEmpty(mevTxs)

	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlow{}
		b := []byte(mevTx.TxFlowPrediction)
		berr := json.Unmarshal(b, &tf)
		s.Require().Nil(berr)
		s.Require().NotEmpty(tf.UserTrade)
		fmt.Println(tf.TradeMethod)
		maxProfit := uni.TradeSimStep(tf)
		fmt.Println("profit from execution path", maxProfit.String())
		//maxProfit := uni.TradeSim(tf)
		//fmt.Println(string(b))
		//fmt.Println("linear search max profit", maxProfit.String())
		fmt.Println("binary search max profit", tf.SandwichPrediction.ExpectedProfit.String())
	}
}
