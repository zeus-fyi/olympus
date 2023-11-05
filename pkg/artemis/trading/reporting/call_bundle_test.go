package artemis_reporting

import (
	"github.com/metachris/flashbotsrpc"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (s *ReportingTestSuite) TestInsertCallBundleResp() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	cr := flashbotsrpc.FlashbotsCallBundleResponse{
		BundleGasPrice:   "",
		BundleHash:       "0x",
		StateBlockNumber: 1,
		TotalGasUsed:     1,
	}
	err := InsertCallBundleResp(ctx, "flashbots", 1, cr)
	s.Assert().Nil(err)
}
