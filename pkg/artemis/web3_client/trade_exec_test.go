package web3_client

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (s *Web3ClientTestSuite) TestTradeExec() {
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
		if tf.FrontRunTrade.AmountIn == "" {
			continue
		}
		_, err := s.LocalHardhatMainnetUser.HardhatResetNetworkToBlockBeforeTxMined(ctx, s.Tc.HardhatNode, s.LocalHardhatMainnetUser, s.MainnetWeb3User, common.HexToHash(tf.Tx.Hash))
		s.Require().Nil(err)

		tfRegular := tf.ConvertToBigIntType()
		err = s.LocalHardhatMainnetUser.MatchFrontRunTradeValues(&tfRegular)
		s.Require().Nil(err)

		uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
		uni.DebugPrint = true

		_, err = uni.ExecFrontRunTradeStepTokenTransfer(&tfRegular)
		s.Require().Nil(err)

		_, _ = uni.UserTradeGetAmountsOut(&tfRegular)
		// must exceed 33073549076721602

		_, err = uni.ExecUserTradeStep(&tfRegular)
		s.Require().Nil(err)

		_, _ = uni.SandwichTradeGetAmountsOut(&tfRegular)
		s.Require().Nil(err)

		_, err = uni.ExecSandwichTradeStepTokenTransfer(&tfRegular)
		s.Require().Nil(err)

		err = tfRegular.GetAggregateGasUsage(ctx, s.LocalHardhatMainnetUser)
		s.Require().Nil(err)

		userGasUsage := tfRegular.UserTrade.TotalGasCost
		fmt.Println("userGasUsage", userGasUsage, "calc", 36627061988*114409)
	}
}
