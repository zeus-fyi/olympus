package web3_client

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
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
	mevTxs, merr := artemis_mev_models.SelectMempoolTxAtMaxBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)

	/*
			should be fully working
			blockNum 17375815
			tradeMethod swapExactETHForTokens
			txHash 0xd40864c0f1d3ad3d2fe4c8e678460d36c4310facfb6be839ca2912c396ef709e

			captures a valid tx with an expired deadline that should be caught before submitting any trades
			blockNum 17375781

		    multi trade
			blockNum 17375834

	*/
	//mevTxs, merr = artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17375869)
	//s.Require().Nil(merr)
	//s.Require().NotEmpty(mevTxs)
	fmt.Println("mevTxs count", len(mevTxs))
	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlowJSON{}
		by := []byte(mevTx.TxFlowPrediction)
		berr := json.Unmarshal(by, &tf)
		s.Assert().Nil(berr)
		if tf.FrontRunTrade.AmountIn == "0" || tf.FrontRunTrade.AmountIn == "" {
			continue
		}
		fmt.Println("blockNum", tf.CurrentBlockNumber)
		fmt.Println("tradeMethod", tf.Trade.TradeMethod)
		fmt.Println("txHash", common.HexToHash(tf.Tx.Hash).Hex())

		rxBlockNum, err := s.LocalHardhatMainnetUser.HardhatResetNetworkToBlockBeforeTxMined(ctx, s.Tc.HardhatNode, s.LocalHardhatMainnetUser, s.MainnetWeb3User, common.HexToHash(tf.Tx.Hash))
		s.Assert().Nil(err)
		blockBeforeRx := rxBlockNum - 1
		tfRegular, err := tf.ConvertToBigIntType()
		s.Assert().Nil(err)
		uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
		uni.DebugPrint = true

		currentBlockStr := tf.CurrentBlockNumber.String()
		currentBlockNum, err := strconv.Atoi(currentBlockStr)
		s.Require().Nil(err)
		fmt.Println("rxBlockNum", rxBlockNum)
		fmt.Println("blockNum recorded from artemis", currentBlockNum)
		fmt.Println("rxBlockNum - artemisBlock", rxBlockNum-currentBlockNum)

		if currentBlockNum < blockBeforeRx {
			fmt.Println("using block number from artemis vs rx block num -1")
			err = s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, currentBlockNum)
			if err != nil {
				log.Err(err).Msg("error resetting hardhat network")
				continue
			}
		}
		if currentBlockNum > rxBlockNum {
			s.Failf("currentBlockNum > rxBlockNum", "currentBlockNum %v > rxBlockNum %v", currentBlockNum, rxBlockNum)
			continue
		}

		//pairAddr := tfRegular.InitialPair.PairContractAddr
		//simPair, err := uni.GetPairContractPrices(ctx, pairAddr)
		//s.Assert().Nil(err)
		//s.Assert().Equal(tfRegular.InitialPair.Reserve0.String(), simPair.Reserve0.String())
		//s.Assert().Equal(tfRegular.InitialPair.Reserve1.String(), simPair.Reserve1.String())

		//err = uni.SimFrontRunTradeOnly(&tfRegular)
		//err = uni.SimUserOnlyTrade(&tfRegular)
		err = uni.SimFullSandwichTrade(&tfRegular)
		s.Assert().Nil(err)
	}
}
