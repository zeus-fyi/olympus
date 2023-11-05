package artemis_reporting

import (
	"fmt"

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

func (s *ReportingTestSuite) TestSelectCallBundles() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	rw, err := SelectCallBundleHistory(ctx, 1699212099729067082, 1)
	s.Assert().Nil(err)

	s.Require().NotNil(rw)

	for _, v := range rw {
		fmt.Println(v.BundleHash)
		fmt.Println(v.Results)
	}
}
