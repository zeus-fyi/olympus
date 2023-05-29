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

	// 17367361
	mevTxs, merr = artemis_validator_service_groups_models.SelectMempoolTxAtBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID, 17367361)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)

	for _, mevTx := range mevTxs {
		tf := TradeExecutionFlow{}
		by := []byte(mevTx.TxFlowPrediction)
		berr := json.Unmarshal(by, &tf)
		s.Assert().Nil(berr)
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

		//err = uni.SimUserOnlyTrade(&tfRegular)

		err = uni.SimFullSandwichTrade(&tfRegular)
		fmt.Println(err)
		err = nil
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

blockNum 17354245
tradeMethod swapExactETHForTokens
txHash 0xdb714f01986223f24dad83b4b358b7be60efc9ef4ac1ba28b2176166f564d2d6
"Error: VM Exception while processing transaction: reverted with reason string 'UniswapV2: INSUFFICIENT_OUTPUT_AMOUNT'

blockNum 17354253
tradeMethod swapETHForExactTokens
txHash 0xa20a8c998191ccc779cbefaad8f51324be2a06a9d81b44c2a2ce9db32eb1a52d
executing full sandwich trade
executing front run trade
executing front run trade
{"level":"error","error":"Error: VM Exception while processing transaction: reverted with reason string 'UniswapV2: INSUFFICIENT_OUTPUT_AMOUNT'","time":"2023-05-27T18:29:28-07:00","message":"error executing front run trade step token transfer"}
Error: VM Exception while processing transaction: reverted with reason string 'UniswapV2: INSUFFICIENT_OUTPUT_AMOUNT'

blockNum 17354304
tradeMethod swapTokensForExactETH
txHash 0xe166c3d418b084f1359c7f9f5f71b6ab0254760d6ad4b10678ed23976a8c8347
executing full sandwich trade
{"level":"error","error":"unable to overwrite balance","time":"2023-05-27T18:40:41-07:00","message":"error executing front run balance setup"}
unable to overwrite balance

*/
