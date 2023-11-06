package artemis_eth_txs

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
)

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
	err := etx.SelectNextPermit2Nonce(ctx)
	s.Require().Nil(err)
	fmt.Println(etx.Nonce)
	fmt.Println(etx.NextPermit2Nonce)

	//err = etx.SelectNextUserTxNonce(ctx, pt)
	//s.Require().Nil(err)
	//fmt.Println(etx.EventID)
	//fmt.Println(etx.NextUserNonce)
}

const (
	AccountAddr = "0x000000641e80A183c8B736141cbE313E136bc8c6"
)

func (s *TxTestSuite) TestSelectExternalUntrackedTxs() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	addrList, err := SelectExternalTxs(ctx, AccountAddr, 1, 0)
	s.Require().Nil(err)
	fmt.Println(addrList)

	//err = etx.SelectNextUserTxNonce(ctx, pt)
	//s.Require().Nil(err)
	//fmt.Println(etx.EventID)
	//fmt.Println(etx.NextUserNonce)
}
