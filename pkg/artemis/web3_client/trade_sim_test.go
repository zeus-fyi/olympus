package web3_client

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

/*
blockNum 17354253
tradeMethod swapETHForExactTokens
txHash 0xa20a8c998191ccc779cbefaad8f51324be2a06a9d81b44c2a2ce9db32eb1a52d
frontRunToken = 1

executing full sandwich trade
executing front run trade
{"level":"error","error":"Error: VM Exception while processing transaction: reverted with reason string 'UniswapV2:
INSUFFICIENT_OUTPUT_AMOUNT'","time":"2023-05-27T18:29:28-07:00","message":"error executing front run trade step token transfer"}
Error: VM Exception while processing transaction: reverted with reason string 'UniswapV2: INSUFFICIENT_OUTPUT_AMOUNT'
*/

/*
blockNum 17354304
tradeMethod swapTokensForExactETH
txHash 0xe166c3d418b084f1359c7f9f5f71b6ab0254760d6ad4b10678ed23976a8c8347
executing full sandwich trade
{"level":"error","error":"unable to overwrite balance","time":"2023-05-27T18:40:41-07:00","message":"error executing front run balance setup"}
unable to overwrite balance
*/

func (s *Web3ClientTestSuite) TestFullSandwichTradeSimAny() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_validator_service_groups_models.SelectMempoolTxAtMaxBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)
	//
	mevTxs, merr = artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17368181)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)

	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlow{}
		by := []byte(mevTx.TxFlowPrediction)
		berr := json.Unmarshal(by, &tf)
		s.Assert().Nil(berr)
		if tf.FrontRunTrade.AmountIn == "0" || tf.FrontRunTrade.AmountIn == "" {
			continue
		}
		fmt.Println("blockNum", tf.CurrentBlockNumber)
		fmt.Println("tradeMethod", tf.Trade.TradeMethod)
		fmt.Println("txHash", tf.Tx.Hash())

		rxBlockNum, err := s.LocalHardhatMainnetUser.HardhatResetNetworkToBlockBeforeTxMined(ctx, s.Tc.HardhatNode, s.LocalHardhatMainnetUser, s.MainnetWeb3User, tf.Tx.Hash())
		s.Assert().Nil(err)
		tfRegular := tf.ConvertToBigIntType()
		uni := InitUniswapV2Client(ctx, s.LocalHardhatMainnetUser)
		uni.DebugPrint = true

		currentBlockStr := tf.CurrentBlockNumber.String()
		currentBlockNum, err := strconv.Atoi(currentBlockStr)
		s.Require().Nil(err)
		fmt.Println("rxBlockNum", rxBlockNum)
		fmt.Println("blockNum recorded from artemis", currentBlockNum)
		fmt.Println("rxBlockNum - artemisBlock", rxBlockNum-currentBlockNum)

		if currentBlockNum < rxBlockNum-1 {
			fmt.Println("using block number from artemis vs rx block num -1")
			err = s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, s.Tc.HardhatNode, currentBlockNum)
			s.Require().Nil(err)
		}
		if currentBlockNum > rxBlockNum {
			s.Failf("currentBlockNum > rxBlockNum", "currentBlockNum %v > rxBlockNum %v", currentBlockNum, rxBlockNum)
		}

		pairAddr := tfRegular.InitialPair.PairContractAddr
		simPair, err := uni.GetPairContractPrices(ctx, pairAddr)
		s.Require().Nil(err)
		s.Require().Equal(tfRegular.InitialPair.Reserve0.String(), simPair.Reserve0.String())
		s.Require().Equal(tfRegular.InitialPair.Reserve1.String(), simPair.Reserve1.String())

		//err = uni.SimFrontRunTradeOnly(&tfRegular)
		//err = uni.SimUserOnlyTrade(&tfRegular)
		//err = uni.SimFullSandwichTrade(&tfRegular)
	}
}
