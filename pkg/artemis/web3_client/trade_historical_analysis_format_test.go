package web3_client

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (s *Web3ClientTestSuite) TestHistoricalAnalysis() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ForceDirToTestDirLocation()
	mevTxs, merr := artemis_validator_service_groups_models.SelectMempoolTxAtMaxBlockNumber(ctx, hestia_req_types.EthereumMainnetProtocolNetworkID)
	s.Require().Nil(merr)
	s.Require().NotEmpty(mevTxs)
	for _, mevTx := range mevTxs {
		uni := InitUniswapV2Client(ctx, s.LocalHardhatMainnetUser)
		uni.DebugPrint = true
		_ = uni.RunHistoricalTradeAnalysis(ctx, mevTx.TxFlowPrediction, s.MainnetWeb3UserExternal)
		uni.PrintResults()
		err := uni.TradeAnalysisReport.SaveResultsInDb(ctx)
		s.Require().Nil(err)
	}
}

/*
Artemis Block Number: 17390664
Rx Block Number: 17390665
End Reason: unable to overwrite balance
End Stage: executing front run balance setup
*/
