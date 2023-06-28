package v1_tyche

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	tyche_base_test "github.com/zeus-fyi/olympus/tyche/api/test"
)

type TxProcessorTestSuite struct {
	tyche_base_test.TycheBaseTestSuite
}

var ctx = context.Background()

func (t *TxProcessorTestSuite) TestTxIngestion() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.LocalDbPgconn)
	t.E.POST(txProcessorRoute, TxProcessingRequestHandler)
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9000")
	}()
	<-start
	defer t.E.Shutdown(ctx)
	txp := TxProcessingRequest{
		Txs: []*types.Transaction{},
	}
	resp, err := t.PostRequest(ctx, txProcessorRoute, txp)
	t.Require().Nil(err)
	fmt.Println(resp)
}

func TestTxProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(TxProcessorTestSuite))
}
