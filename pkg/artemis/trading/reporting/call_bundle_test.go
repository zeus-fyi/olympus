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
		BundleHash:       "0x1",
		StateBlockNumber: 1,
		Results: []flashbotsrpc.FlashbotsCallBundleResult{
			{
				CoinbaseDiff:      "",
				EthSentToCoinbase: "",
				FromAddress:       "",
				GasFees:           "",
				GasPrice:          "",
				GasUsed:           0,
				ToAddress:         "",
				TxHash:            "",
				Value:             "",
				Error:             "\u0000",
			},
		},
		TotalGasUsed: 1,
	}

	//out := clearString("\u0008�y�\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000 \u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0014TRANSFER_FROM_FAILED\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000")
	//fmt.Println(out)
	err := InsertCallBundleResp(ctx, "flashbots", 1, cr, nil)
	s.Assert().Nil(err)
}

func (s *ReportingTestSuite) TestSelectCallBundles() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	rw, err := SelectCallBundleHistory(ctx, 0, 1)
	s.Assert().Nil(err)

	s.Require().NotNil(rw)

	for _, v := range rw {
		fmt.Println("=====================================")
		fmt.Println(v.BundleHash)
		fmt.Println(v.BundleGasPrice)
		fmt.Println(v.TotalGasUsed)

		for _, res := range v.Results {
			if len(res.Error) > 0 {
				fmt.Println("Error", res.Error)
			}
			if len(res.Revert) > 0 {
				fmt.Println("Revert", res.Revert)
			}
			fmt.Println(res.GasUsed)
			fmt.Println(res.GasPrice)
			fmt.Println(res.GasFees)
		}
		fmt.Println("=====================================")
	}
}
