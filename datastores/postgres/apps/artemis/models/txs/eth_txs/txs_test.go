package artemis_eth_txs

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
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
	s.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	etx := EthTx{
		EthTxGas: artemis_autogen_bases.EthTxGas{
			TxHash:   "0x000w25e60C7ff32a3470be7FE3ed1666b0E326e2",
			GasPrice: sql.NullInt64{},
			GasLimit: sql.NullInt64{
				Int64: 300000,
			},
			GasTipCap: sql.NullInt64{
				Int64: artemis_eth_units.GweiMultiple(2).Int64(),
			},
			GasFeeCap: sql.NullInt64{
				Int64: artemis_eth_units.GweiMultiple(30).Int64(),
			},
		},
		EthTx: artemis_autogen_bases.EthTx{
			ProtocolNetworkID: hestia_req_types.EthereumGoerliProtocolNetworkID,
			TxHash:            "0x0ee12f11d",
			Nonce:             0,
			From:              "0x000025e60C7ff32a3470be7FE3ed1666b0E326e2",
			Type:              "0x02",
			EventID:           5,
		},
	}
	pt := Permit2Tx{Permit2Tx: artemis_autogen_bases.Permit2Tx{
		Nonce:    0,
		Owner:    "0x000025e60C7ff32a3470be7FE3ed1666b0E326e2",
		Deadline: int(time.Now().Add(time.Minute * 5).Unix()),
		EventID:  int(time.Now().Unix()),
		Token:    artemis_trading_constants.GoerliWETH9ContractAddress,
	}}
	err := etx.InsertTx(ctx, pt)
	s.Require().Nil(err)
}

func (s *TxTestSuite) TestInsertTx() {
	etx := EthTx{
		EthTx: artemis_autogen_bases.EthTx{
			ProtocolNetworkID: 5,
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

func (s *TxTestSuite) TestInsertBundle() {
	pi := hestia_req_types.EthereumGoerliProtocolNetworkID

	bundleTxs := []EthTx{
		{
			EthTx: artemis_autogen_bases.EthTx{
				ProtocolNetworkID: pi,
				TxHash:            "0x012fad",
				Nonce:             0,
				From:              "0x0gsdg32",
				Type:              "0x02",
			},
			EthTxGas: artemis_autogen_bases.EthTxGas{
				TxHash:    "",
				GasPrice:  sql.NullInt64{},
				GasLimit:  sql.NullInt64{},
				GasTipCap: sql.NullInt64{},
				GasFeeCap: sql.NullInt64{},
			},
			Permit2Tx: Permit2Tx{
				Permit2Tx: artemis_autogen_bases.Permit2Tx{
					Nonce:             0,
					Owner:             "",
					Deadline:          0,
					Token:             "",
					ProtocolNetworkID: pi,
				},
			},
		},
		{
			EthTx: artemis_autogen_bases.EthTx{
				ProtocolNetworkID: pi,
				TxHash:            "",
				Nonce:             0,
				From:              "",
				Type:              "",
			},
			EthTxGas: artemis_autogen_bases.EthTxGas{
				TxHash:    "",
				GasPrice:  sql.NullInt64{},
				GasLimit:  sql.NullInt64{},
				GasTipCap: sql.NullInt64{},
				GasFeeCap: sql.NullInt64{},
			},
		},
		{
			EthTx: artemis_autogen_bases.EthTx{
				ProtocolNetworkID: pi,
			},
			EthTxGas: artemis_autogen_bases.EthTxGas{
				TxHash:    "",
				GasPrice:  sql.NullInt64{},
				GasLimit:  sql.NullInt64{},
				GasTipCap: sql.NullInt64{},
				GasFeeCap: sql.NullInt64{},
			},
			Permit2Tx: Permit2Tx{
				Permit2Tx: artemis_autogen_bases.Permit2Tx{
					Nonce:             0,
					Owner:             "",
					Deadline:          0,
					Token:             "0xs",
					ProtocolNetworkID: pi,
				},
			},
		},
	}
	bundleHash := "0x012fad"
	err := InsertTxsWithBundle(ctx, bundleTxs, bundleHash)
	s.Require().Nil(err)

}

func TestTxTestSuite(t *testing.T) {
	suite.Run(t, new(TxTestSuite))
}
