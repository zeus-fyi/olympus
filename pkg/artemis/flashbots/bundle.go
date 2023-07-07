package artemis_flashbots

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/metachris/flashbotsrpc"
)

type MevTxBundle struct {
	flashbotsrpc.FlashbotsSendBundleRequest
}

func (m *MevTxBundle) AddTxs(txs ...*types.Transaction) {
	if m.Txs == nil {
		m.Txs = []string{}
	}
	for _, tx := range txs {
		m.Txs = append(m.Txs, tx.Hash().String())
	}
}
