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
	mevTxs, merr := artemis_validator_service_groups_models.SelectEthMevTxAnalysisByTxHash(ctx, "0xb653a3cf525a0fc9df716203d893fe1479422aed29489efd0acfd31bb09dec9b")
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)
	for _, mevTx := range mevTxs {
		fmt.Println(mevTx.TradeMethod)
		uni := InitUniswapV2Client(ctx, s.HostedHardhatMainnetUser)
		uni.DebugPrint = true
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
		uni := InitUniswapV2Client(ctx, s.HostedHardhatMainnetUser)
		uni.DebugPrint = true
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
