package artemis_eth_txs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

var ctx = context.Background()

type TxTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *TxTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
}

func (s *TxTestSuite) TestInsertPermit2Tx() {
	etx := EthTx{
		artemis_autogen_bases.EthTx{
			ProtocolNetworkID: 0,
			TxHash:            "",
			Nonce:             0,
			From:              "",
			Type:              "",
			EventID:           0,
		},
	}
	err := etx.InsertTx(ctx)
	s.Require().Nil(err)
}

func (s *TxTestSuite) TestInsertTx() {
	etx := EthTx{
		artemis_autogen_bases.EthTx{
			ProtocolNetworkID: 1,
			TxHash:            "0x012fad",
			Nonce:             0,
			From:              "0x0gsdg32",
			Type:              "0x02",
			EventID:           0,
		},
	}
	err := etx.InsertTx(ctx)
	s.Require().Nil(err)
	s.Assert().NotZerof(etx.EventID, "event id should not be zero")
}

func (s *TxTestSuite) TestSelect() {
}

func TestTxTestSuite(t *testing.T) {
	suite.Run(t, new(TxTestSuite))
}
