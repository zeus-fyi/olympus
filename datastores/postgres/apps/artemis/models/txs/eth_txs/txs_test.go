package artemis_eth_txs

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
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
		EthTxGas: artemis_autogen_bases.EthTxGas{
			TxHash:   "0x0ee12f11d",
			GasPrice: sql.NullInt64{},
			GasLimit: sql.NullInt64{
				Int64: 300000,
			},
			GasTipCap: sql.NullInt64{
				Int64: artemis_eth_units.GweiMultiple(15).Int64(),
			},
			GasFeeCap: sql.NullInt64{
				Int64: artemis_eth_units.GweiMultiple(1).Int64(),
			},
		},
		EthTx: artemis_autogen_bases.EthTx{
			ProtocolNetworkID: 1,
			TxHash:            "0x0ee12f11d",
			Nonce:             1,
			From:              "0x0gsdg32",
			Type:              "0x02",
			EventID:           1,
		},
	}
	pt := Permit2Tx{Permit2Tx: artemis_autogen_bases.Permit2Tx{
		Nonce:    1,
		Owner:    "0x0gsdg32",
		Deadline: int(time.Now().Add(time.Minute * 5).Unix()),
		EventID:  int(time.Now().Unix()),
		Token:    "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
	}}
	err := etx.InsertTx(ctx, pt)
	s.Require().Nil(err)
}

func (s *TxTestSuite) TestInsertTx() {
	etx := EthTx{
		EthTx: artemis_autogen_bases.EthTx{
			ProtocolNetworkID: 1,
			TxHash:            "0x012fad",
			Nonce:             0,
			From:              "0x0gsdg32",
			Type:              "0x02",
			EventID:           0,
		},
	}
	//err := etx.InsertTx(ctx)
	//s.Require().Nil(err)
	s.Assert().NotZerof(etx.EventID, "event id should not be zero")
}

func (s *TxTestSuite) TestSelect() {
	etx := EthTx{
		EthTx: artemis_autogen_bases.EthTx{
			ProtocolNetworkID: 1,
			TxHash:            "0x012fad",
			Nonce:             0,
			From:              "0x0gsdg32",
			Type:              "0x02",
			EventID:           0,
		},
	}
	pt := Permit2Tx{Permit2Tx: artemis_autogen_bases.Permit2Tx{
		Nonce:    1,
		Owner:    "0x0gsdg32",
		Deadline: int(time.Now().Add(time.Minute * 5).Unix()),
		EventID:  int(time.Now().Unix()),
		Token:    "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
	}}
	err := etx.SelectNextPermit2Nonce(ctx, pt)
	s.Require().Nil(err)
	fmt.Println(etx.Nonce)
	fmt.Println(etx.NextPermit2Nonce)

	//err = etx.SelectNextUserTxNonce(ctx, pt)
	//s.Require().Nil(err)
	//fmt.Println(etx.EventID)
	//fmt.Println(etx.NextUserNonce)
}

func TestTxTestSuite(t *testing.T) {
	suite.Run(t, new(TxTestSuite))
}
