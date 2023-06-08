package web3_client

import (
	"context"
	"math/big"
)

const (
	SudoSwap = "SUDOSWAP"
)

type SudoSwapParams struct {
	Value *big.Int `json:"value"`
	Data  []byte   `json:"bytes"`
}

func (s *SudoSwapParams) Encode(ctx context.Context) ([]byte, error) {
	return nil, nil
}

func (s *SudoSwapParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoder.Methods[SudoSwap].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	value, err := ParseBigInt(args["value"])
	if err != nil {
		return err
	}
	s.Value = value
	s.Data = args["data"].([]byte)
	return nil
}
