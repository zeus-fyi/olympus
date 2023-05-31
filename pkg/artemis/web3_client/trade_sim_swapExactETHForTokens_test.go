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
	blockNum 17375869
	tradeMethod swapExactETHForTokens
	txHash 0xe11f91fe084c02eb92b72c42a053edade40a28c62b459f3404d7179face2e7f5
*/

func (s *Web3ClientTestSuite) TestFullSandwichTradeSim_SwapExactETHForTokens() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17375869)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)
	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlow{}
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
		err = s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, s.Tc.HardhatNode, currentBlockNum)
		s.Require().Nil(err)

		tfRegular := tf.ConvertToBigIntType()
		uni := InitUniswapV2Client(ctx, s.LocalHardhatMainnetUser)
		pairAddr := tfRegular.InitialPair.PairContractAddr
		simPair, err := uni.GetPairContractPrices(ctx, pairAddr)
		s.Require().Nil(err)
		s.Require().Equal(tfRegular.InitialPair.Reserve0.String(), simPair.Reserve0.String())
		s.Require().Equal(tfRegular.InitialPair.Reserve1.String(), simPair.Reserve1.String())

		uni.DebugPrint = true
		err = uni.SimFullSandwichTrade(&tfRegular)
		s.Require().Nil(err)
	}
}
