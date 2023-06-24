package web3_client

import (
	"context"
	"math/big"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

type Multicall struct {
	Deadline *big.Int `json:"deadline,omitempty"`
	Data     [][]byte `json:"data"`
}

func (m *Multicall) Decode(ctx context.Context, args map[string]interface{}) error {
	deadline, ok := args["deadline"].(*big.Int)
	if !ok {
		log.Info().Msg("failed to decode deadline")
	} else {
		m.Deadline = deadline
	}
	data, ok := args["data"].([][]byte)
	if !ok {
		return errors.New("failed to decode data")
	}
	m.Data = data
	return nil
}
