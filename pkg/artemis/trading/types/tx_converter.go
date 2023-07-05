package artemis_trading_types

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"
)

type JSONTx struct {
	Type                 string      `json:"type"`
	Nonce                string      `json:"nonce"`
	To                   string      `json:"to"`
	Gas                  string      `json:"gas"`
	ChainID              string      `json:"chainId"`
	GasPrice             string      `json:"gasPrice,omitempty"`
	MaxPriorityFeePerGas interface{} `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerGas         interface{} `json:"maxFeePerGas,omitempty"`
	Value                string      `json:"value,omitempty"`
	Input                string      `json:"input"`
	V                    string      `json:"v"`
	R                    string      `json:"r"`
	S                    string      `json:"s"`
	Hash                 string      `json:"hash"`
}

func (j *JSONTx) UnmarshalTx(tx *types.Transaction) error {
	b, err := json.Marshal(tx)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	return nil
}

func (j *JSONTx) ConvertToTx() (*types.Transaction, error) {
	tx := &types.Transaction{}
	b, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
