package artemis_eth_txs

import (
	"fmt"

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
