package web3_client

import (
	"encoding/json"
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (s *Web3ClientTestSuite) TestFullSandwichTradeSim() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17332397)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)

	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlow{}
		by := []byte(mevTx.TxFlowPrediction)
		berr := json.Unmarshal(by, &tf)
		s.Require().Nil(berr)
		if tf.FrontRunTrade.AmountIn == "" {
			continue
		}
		err := s.LocalHardhatMainnetUser.HardhatResetNetworkToBlockBeforeTxMined(ctx, s.Tc.HardhatNode, s.LocalHardhatMainnetUser, s.MainnetWeb3User, tf.Tx.Hash())
		s.Require().Nil(err)
		tfRegular := tf.ConvertToBigIntType()
		uni := InitUniswapV2Client(ctx, s.LocalHardhatMainnetUser)
		uni.DebugPrint = true
		err = uni.SimFullSandwichTrade(&tfRegular)
		s.Require().Nil(err)
	}
}

func (s *Web3ClientTestSuite) TestFullSandwichTradeSimAny() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_validator_service_groups_models.SelectMempoolTxAtMaxBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)

	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlow{}
		by := []byte(mevTx.TxFlowPrediction)
		berr := json.Unmarshal(by, &tf)
		s.Require().Nil(berr)
		if tf.FrontRunTrade.AmountIn == "" {
			continue
		}
		fmt.Println("blockNum", tf.CurrentBlockNumber)
		fmt.Println("tradeMethod", tf.Trade.TradeMethod)
		fmt.Println("txHash", tf.Tx.Hash())

		err := s.LocalHardhatMainnetUser.HardhatResetNetworkToBlockBeforeTxMined(ctx, s.Tc.HardhatNode, s.LocalHardhatMainnetUser, s.MainnetWeb3User, tf.Tx.Hash())
		s.Assert().Nil(err)
		tfRegular := tf.ConvertToBigIntType()
		uni := InitUniswapV2Client(ctx, s.LocalHardhatMainnetUser)
		uni.DebugPrint = true
		err = uni.SimFullSandwichTrade(&tfRegular)
		s.Assert().Nil(err)
	}
}

/*
blockNum 17354228
TRANSFER_FAILED
tradeMethod swapTokensForExactTokens
txHash 0x5f29de88cd07de0582923590cdaf77dcc35c73d458549e6ec103f8e5e80b06ed

blockNum 17354235
tradeMethod swapETHForExactTokens
txHash 0xfc3ae1c4ef163d8a974bea83dd23f7a81c168dc02c61a2d4d5f536223d683509
VM Exception while processing transaction: reverted with reason string 'SafeMath: subtraction overflow'
*/
