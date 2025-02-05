package web3_client

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

/*
blockNum 17332397
tradeMethod swapExactTokensForETH
txHash 0xed1d8212026203c2ae1c84d132f81e0c107f738d53462459fc1b68cd0f97b743
frontRunTradeToken = 0

blockNum 17384016
tradeMethod swapExactTokensForETH
txHash 0xde078377d909ad0cc5c05e7f854d168ba9192fc41815eaf83f730b0b337eaaac
*/

func (s *Web3ClientTestSuite) TestFullSandwichTradeSim_SwapExactTokensForETH() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_mev_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17332397)
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
		s.Require().Equal(swapExactTokensForETH, tf.Trade.TradeMethod)

		fmt.Println("blockNum recorded from artemis", tf.CurrentBlockNumber)
		fmt.Println("tradeMethod", tf.Trade.TradeMethod)
		fmt.Println("txHash", common.HexToHash(tf.Tx.Hash))
		blockNum, err := s.LocalHardhatMainnetUser.HardhatResetNetworkToBlockBeforeTxMined(ctx, s.Tc.HardhatNode, s.LocalHardhatMainnetUser, s.MainnetWeb3User, common.HexToHash(tf.Tx.Hash))
		fmt.Println("blockNumSet to -1 before tx included", blockNum-1)
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

		currentBlockStr := tf.CurrentBlockNumber.String()
		currentBlockNum, err := strconv.Atoi(currentBlockStr)
		s.Require().Nil(err)
		fmt.Println("blockNum recorded from artemis", currentBlockNum)
		err = s.LocalHardhatMainnetUser.HardHatResetNetwork(ctx, currentBlockNum)
		s.Require().Nil(err)

		tfRegular, err = tf.ConvertToBigIntType()
		s.Require().Nil(err)
		uni = InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
		uni.PrintDetails = true
		uni.TestMode = true

		//pairAddr = tfRegular.InitialPair.PairContractAddr
		//simPair, err = uni.GetPairContractPrices(ctx, pairAddr)
		//s.Require().Nil(err)
		//s.Require().Equal(tfRegular.InitialPair.Reserve0.String(), simPair.Reserve0.String())
		//s.Require().Equal(tfRegular.InitialPair.Reserve1.String(), simPair.Reserve1.String())
		uni.PrintOn = true
		uni.DebugPrint = true
		err = uni.SimFullSandwichTrade(&tfRegular)
		s.Require().Nil(err)
	}
}
