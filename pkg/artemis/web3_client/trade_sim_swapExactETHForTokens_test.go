package web3_client

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

/*
	blockNum 17375869
	tradeMethod swapExactETHForTokens
	txHash 0xe11f91fe084c02eb92b72c42a053edade40a28c62b459f3404d7179face2e7f5

	blockNum 17383929
	tradeMethod swapExactETHForTokens
	txHash 0xeac16063a5968c7c338d869b9283a95b7b1482e7e0cdf687ef98c245c3e54915

	blockNum 17383953
	tradeMethod swapExactETHForTokens
	txHash 0xa3b79b88b70d734aa55d89dd75285b73b56112feca95caad1c861e1583ea4923

	blockNum 17384011
	tradeMethod swapExactETHForTokens
	txHash 0x6762d5ec93238c65421b04f966c7c802caab5d979e62c0d6d090eaf25acf2e12

ERRORS

blockNum 17383946
tradeMethod swapExactETHForTokens
txHash 0x0e305555d8ed6afd7e63fad455a03830a1c3f8ad1c064b77786ec9b2141181a3
{"level":"warn","transferTx":null,"error":"Error: VM Exception while processing transaction: reverted with reason string 'Insufficient Balance'","time":"2023-05-31T22:41:39-07:00","message":"error approving router"}

*/

func (s *Web3ClientTestSuite) TestFullSandwichTradeSim_SwapExactETHForTokens() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_mev_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17375869)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)
	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlowJSON{}
		by := []byte(mevTx.TxFlowPrediction)
		berr := json.Unmarshal(by, &tf)
		s.Require().Nil(berr)
		if tf.FrontRunTrade.AmountIn == "0" || tf.FrontRunTrade.AmountIn == "" {
			continue
		}
		s.Require().Equal(swapExactETHForTokens, tf.Trade.TradeMethod)
		fmt.Println("blockNum recorded from artemis", tf.CurrentBlockNumber)
		currentBlockStr := tf.CurrentBlockNumber.String()
		currentBlockNum, err := strconv.Atoi(currentBlockStr)
		s.Require().Nil(err)
		fmt.Println("blockNum recorded from artemis", currentBlockNum)
		err = s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, currentBlockNum)
		s.Require().Nil(err)

		tfRegular := tf.ConvertToBigIntType()
		uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
		//pairAddr := tfRegular.InitialPair.PairContractAddr
		//simPair, err := uni.GetPairContractPrices(ctx, pairAddr)
		//s.Require().Nil(err)
		//s.Require().Equal(tfRegular.InitialPair.Reserve0.String(), simPair.Reserve0.String())
		//s.Require().Equal(tfRegular.InitialPair.Reserve1.String(), simPair.Reserve1.String())

		uni.DebugPrint = true
		uni.TestMode = true
		err = uni.SimFullSandwichTrade(&tfRegular)
		s.Require().Nil(err)
	}
}
