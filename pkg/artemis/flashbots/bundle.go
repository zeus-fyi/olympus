package artemis_flashbots

import "github.com/ethereum/go-ethereum/core/types"

type MevTxBundle struct {
	Txs []*types.Transaction `json:"txs"`
}
