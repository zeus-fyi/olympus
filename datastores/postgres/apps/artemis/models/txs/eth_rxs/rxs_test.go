package artemis_eth_rxs

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

var ctx = context.Background()

type RxTestSuite struct {
	w3c web3_client.Web3Client
	hestia_test.BaseHestiaTestSuite
}

func (s *RxTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	s.w3c = web3_client.NewWeb3ClientFakeSigner("https://iris.zeus.fyi/v1/router")
	s.w3c.AddDefaultEthereumMainnetTableHeader()
	s.w3c.AddBearerToken(s.Tc.ProductionLocalBearerToken)
}

func (s *RxTestSuite) TestInsertRx() {
	rx, found, rerr := s.w3c.GetTxReceipt(ctx, common.HexToHash("0x"))
	s.Assert().Nil(rerr)
	if !found {
		s.Assert().Fail("tx not found")
	}
	status := "unknown"
	if rx.Status == types.ReceiptStatusSuccessful {
		status = "success"
	}
	if rx.Status == types.ReceiptStatusFailed {
		status = "failed"
	}
	rxEthTx := artemis_autogen_bases.EthTxReceipts{
		Status:            status,
		GasUsed:           int(rx.GasUsed),
		CumulativeGasUsed: int(rx.CumulativeGasUsed),
		BlockHash:         rx.BlockHash.String(),
		TransactionIndex:  int(rx.TransactionIndex),
		TxHash:            rx.TxHash.String(),
		EventID:           0,
		EffectiveGasPrice: int(rx.EffectiveGasPrice.Int64()),
		BlockNumber:       int(rx.BlockNumber.Int64()),
	}
	err := InsertTxReceipt(ctx, rxEthTx)
	s.Assert().Nil(err)
}

/*
status := "unknown"

	if rx.Status == types.ReceiptStatusSuccessful {
		status = "success"
	}
	if rx.Status == types.ReceiptStatusFailed {
		status = "failed"
	}
*/
func TestRxTestSuite(t *testing.T) {
	suite.Run(t, new(RxTestSuite))
}
