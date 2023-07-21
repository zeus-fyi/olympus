package artemis_realtime_trading

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

var ctx = context.Background()

func (s *ArtemisRealTimeTradingTestSuite) TestPipeline() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	mevMempoolTx, err := artemis_mev_models.SelectEthMevMempoolTxByTxHash(ctx, "0xb329c536560589e2ee14f948cd84e6f05d59e9b18b95dba257b618a9bb3d16b0")
	s.Require().NoError(err)

	for _, mevTx := range mevMempoolTx {
		fmt.Println("tradeMethod", mevTx.TradeMethod)

		strTx := mevTx.EthMempoolMevTx.Tx
		jtx := artemis_trading_types.JSONTx{}
		err = json.Unmarshal([]byte(strTx), &jtx)
		s.Require().NoError(err)

		tx, terr := jtx.ConvertToTx()
		s.Require().NoError(terr)

		fmt.Println("tx", tx.Hash().String())
		fmt.Println("txHash", mevTx.EthMempoolMevTx.Tx)
	}
	//IngestTx(ctx

}
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
