package web3_client

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (s *Web3ClientTestSuite) TestHistoricalAnalysisReplay() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_mev_models.SelectEthMevTxAnalysisByTxHash(ctx, "0xf47b8eb14f06db38ab6e23f04f86e23ae0c0797d62a1f2780e5bba93343a666d")
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
	mevTxs, merr := artemis_mev_models.SelectMempoolTxAtMaxBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID)
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
