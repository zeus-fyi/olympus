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
	tmp, err := GetHexEncodedTxStrSlice(txs...)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tmp...)
	return nil
}

func GetHexEncodedTxStrSlice(txs ...*types.Transaction) ([]string, error) {
	txSlice := make([]string, len(txs))
	for i, tx := range txs {
		b, err := tx.MarshalBinary()
		if err != nil {
			return nil, err
		}
		txSlice[i] = hexutil.Encode(b)
	}
	return txSlice, nil
}
