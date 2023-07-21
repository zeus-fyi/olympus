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
	mevTxs count 1
	blockNum 17383844
	tradeMethod swapExactTokensForTokens
	txHash 0x96837858590d7805a3503ddfadd5f4dbb3fdad3a7291234c0e3bafe0bb2261d4
*/

func (s *Web3ClientTestSuite) TestFullSandwichTradeSim_SwapExactTokensForTokens() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_mev_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17383844)
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
		s.Require().Equal(swapExactTokensForTokens, tf.Trade.TradeMethod)
		fmt.Println("blockNum recorded from artemis", tf.CurrentBlockNumber)
		currentBlockStr := tf.CurrentBlockNumber.String()
		currentBlockNum, err := strconv.Atoi(currentBlockStr)
		s.Require().Nil(err)
		fmt.Println("blockNum recorded from artemis", currentBlockNum)
		err = s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, currentBlockNum)
		s.Require().Nil(err)

		tfRegular, err := tf.ConvertToBigIntType()
		s.Require().Nil(err)
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
