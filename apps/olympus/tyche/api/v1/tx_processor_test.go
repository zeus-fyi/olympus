package v1_tyche

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

	hashStr := "0xb841ae58afb7c6e0e7c321e2d151d93599dfd826ac3835f3c7cd8c029b6fd9a7"
	tx, _, err := t.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	t.Require().Nil(err)
	t.Require().NotNil(tx)

	txp := TxProcessingRequest{
		Txs: []*types.Transaction{tx},
	}
	resp, err := t.PostRequest(ctx, txProcessorRoute, txp)
	t.Require().Nil(err)
	fmt.Println(resp)
}

func TestTxProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(TxProcessorTestSuite))
}
