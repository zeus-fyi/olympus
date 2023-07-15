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
	mevMempoolTx, err := artemis_mev_models.SelectEthMevMempoolTxByTxHash(ctx, "0xb647f1939f1ebd9a9de5cd9fe2ae869cd7f8bfc1d34879e29a0ccbdec96876c4")
	s.Require().NoError(err)
	s.Require().Len(mevMempoolTx, 1)
	tx := mevMempoolTx[0]
	j, err := web3_client.UnmarshalTradeExecutionFlow(tx.TxFlowPrediction)
	s.Require().NoError(err)

	tmp := j.ConvertToBigIntType()
	s.Require().NotEmpty(tmp)
	fmt.Println(j.ConvertToBigIntType())

	err = s.at.SaveMempoolTx(ctx, tmp.CurrentBlockNumber.Uint64(), []web3_client.TradeExecutionFlowJSON{j})
	s.Require().NoError(err)

}
