package web3_client

import (
	"math/big"
)

type UniversalRouterExecCmd struct {
	Commands []UniversalRouterExecSubCmd `json:"commands"`
	Deadline *big.Int                    `json:"deadline"`
}

type UniversalRouterExecSubCmd struct {
	Command       string `json:"command"`
	CanRevert     bool   `json:"canRevert"`
	Inputs        []byte `json:"inputs"`
	DecodedInputs any    `json:"decodedInputs"`
}
