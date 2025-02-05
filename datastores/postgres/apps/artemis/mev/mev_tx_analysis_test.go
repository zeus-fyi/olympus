package artemis_mev_models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type MevTxAnalysisTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *MevTxAnalysisTestSuite) TestSelectTradeMethodStatsByPair() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	tms, err := SelectEthMevTxAnalysisByPair(ctx)
	s.Require().Nil(err)
	s.Require().NotEmpty(tms)
	for _, tm := range tms {
		fmt.Println(tm.AmountInAddr, tm.AmountOutAddr)
	}
}

func (s *MevTxAnalysisTestSuite) TestSelectTradeMethodStatsBySuccess() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	tms, err := SelectTradeMethodStatsBySuccess(ctx)
	s.Require().Nil(err)
	s.Require().NotEmpty(tms)
	for _, tm := range tms {
		fmt.Println(tm.TradeMethod, tm.Count)
	}
}

func (s *MevTxAnalysisTestSuite) TestInsertNodes() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	txHistory := artemis_autogen_bases.EthMevTxAnalysis{
		GasUsedWei:              "1000",
		Metadata:                "{}",
		TxHash:                  "0x5a26a6207b24770ee69e82c318163136fc4d96758d68a56bc63efae25f6a394d",
		TradeMethod:             "swap",
		EndReason:               "success",
		AmountIn:                "0",
		AmountOutAddr:           "0x",
		ExpectedProfitAmountOut: "0",
		RxBlockNumber:           100,
		AmountInAddr:            "0x",
		ActualProfitAmountOut:   "0",
	}

	err := InsertEthMevTxAnalysis(ctx, txHistory)
	s.Require().Nil(err)
}

func (s *MevTxAnalysisTestSuite) TestSelectTxAnalysis() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	resp, err := SelectEthMevTxAnalysis(ctx)
	s.Require().Nil(err)
	s.Require().NotEmpty(resp)
}

func TestMevTxAnalysisTestSuite(t *testing.T) {
	suite.Run(t, new(MevTxAnalysisTestSuite))
}
