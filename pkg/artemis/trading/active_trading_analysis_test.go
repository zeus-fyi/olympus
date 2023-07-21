package artemis_realtime_trading

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

var ctx = context.Background()

func (s *ArtemisRealTimeTradingTestSuite) TestCalculateTradeValues() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	mevMempoolTx, err := artemis_mev_models.SelectEthMevMempoolTxByTxHash(ctx, "0x4a94d6c07a8d97bd94f3d940136860f850df41494fd0976e412898313e33bf49")
	s.Require().NoError(err)
	s.Require().Len(mevMempoolTx, 1)
	tx := mevMempoolTx[0]
	j, err := web3_client.UnmarshalTradeExecutionFlow(tx.TxFlowPrediction)
	s.Require().NoError(err)

	tmp, err := j.ConvertToBigIntType()
	s.Require().NoError(err)
	s.Require().NotEmpty(tmp)
	fmt.Println(j.ConvertToBigIntType())

	tp, err := j.Trade.JSONV2SwapExactInParams.ConvertToBigIntType()
	s.Require().NoError(err)
	s.Require().NotEmpty(tp)

	to, err := tp.BinarySearch(*tmp.InitialPair)
	s.Require().NoError(err)

	s.Require().NotEmpty(to)
	// 13802675169811703
	fmt.Println("orig", tmp.SandwichPrediction.ExpectedProfit)
	fmt.Println("sand", to.SandwichPrediction.ExpectedProfit)
}
