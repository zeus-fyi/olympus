package artemis_flashbots

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/metachris/flashbotsrpc"
)

type MevTxBundle struct {
	EventID int
	*flashbotsrpc.FlashbotsSendBundleRequest
}

func (m *MevTxBundle) AddTxs(txs ...*types.Transaction) error {
	if m.Txs == nil {
		m.Txs = []string{}
	}
	tmp := make([]string, len(txs))
	for i, tx := range txs {
		b, err := tx.MarshalBinary()
		if err != nil {
			return err
		}
		tmp[i] = hexutil.Encode(b)
	}
	m.Txs = append(m.Txs, tmp...)
	return nil
}
