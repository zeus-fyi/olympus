package web3_client

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (s *Web3ClientTestSuite) TestHistoricalAnalysisReplay() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_validator_service_groups_models.SelectEthMevTxAnalysisByTxHash(ctx, "0xa864c448e3732160c6aabd4e3e90aad990e07043f1a45bcffe74800cf9d58aff")
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)
	for _, mevTx := range mevTxs {
		fmt.Println(mevTx.TradeMethod)
		uni := InitUniswapClient(ctx, s.ProxyHostedHardhatMainnetUser)
		uni.Web3Client.IsAnvilNode = true
		uni.DebugPrint = true
		uni.PrintLocal = true
		uni.PrintDetails = true
		err := uni.RunHistoricalTradeAnalysis(ctx, mevTx.TxFlowPrediction, s.MainnetWeb3UserExternal)
		uni.PrintResults()
		s.Assert().Nil(err)
	}
}

func (s *Web3ClientTestSuite) TestHistoricalAnalysis() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_validator_service_groups_models.SelectMempoolTxAtMaxBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)
	for _, mevTx := range mevTxs {
		uni := InitUniswapClient(ctx, s.HostedHardhatMainnetUser)
		uni.DebugPrint = true
		uni.PrintLocal = true
		uni.PrintDetails = true
		uni.PrintOn = true
		uni.TestMode = true
		err := uni.RunHistoricalTradeAnalysis(ctx, mevTx.TxFlowPrediction, s.MainnetWeb3UserExternal)
		uni.PrintResults()
		s.Assert().Nil(err)
	}
}

/*
Artemis Block Number: 17390664
Rx Block Number: 17390665
End Reason: unable to overwrite balance
End Stage: executing front run balance setup
*/
