package web3_client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
)

// todo exec
/* priority list
0x00	V3_SWAP_EXACT_IN
0x01	V3_SWAP_EXACT_OUT
0x02	PERMIT2_TRANSFER_FROM
0x03	PERMIT2_PERMIT_BATCH
0x08	V2_SWAP_EXACT_IN
0x09	V2_SWAP_EXACT_OUT
0x0a	PERMIT2_PERMIT
0x0d	PERMIT2_TRANSFER_FROM_BATCH
*/

const (
	SudoSwap = "SUDOSWAP"
)

func (u *UniswapClient) DecodeUniversalRouterMessage() {
	// TODO
	// get command from bytes
}

func NewDecodedUniversalRouterExecCmdFromMap(m map[string]interface{}) (UniversalRouterExecCmd, error) {
	cmds := UniversalRouterExecCmd{
		Commands: []UniversalRouterExecSubCmd{},
	}
	var commandsVal []byte
	var inputsVal [][]byte
	for k, v := range m {
		switch k {
		case "commands":
			commandsVal = v.([]byte)
		case "inputs":
			inputsVal = v.([][]byte)
		case "deadline":
			cmds.Deadline = v.(*big.Int)
		}
	}
	subCmds := make([]UniversalRouterExecSubCmd, len(commandsVal))
	for i, byteSize := range commandsVal {
		subCmd := UniversalRouterExecSubCmd{}
		err := subCmd.DecodeCommand(byteSize)
		if err != nil {
			return cmds, err
		}
		subCmd.Inputs = inputsVal[i]
		subCmds[i] = subCmd
	}
	cmds.Commands = subCmds
	return cmds, nil
}

type UniversalRouterExecCmd struct {
	Commands []UniversalRouterExecSubCmd `json:"commands"`
	Deadline *big.Int                    `json:"deadline"`
}

type UniversalRouterExecSubCmd struct {
	Command   string `json:"command"`
	CanRevert bool   `json:"canRevert"`
	Inputs    []byte `json:"inputs"`
}

func (ur *UniversalRouterExecSubCmd) DecodeCommand(command byte) error {
	buf := bytes.NewBuffer([]byte{command})
	var data uint8
	err := binary.Read(buf, binary.BigEndian, &data)
	if err != nil {
		return fmt.Errorf("could not read command: %v", err)
	}
	flag := (data & 0x80) >> 7 // extract bit 7
	//ref := (data & 0x60) >> 5  // extract bits 6-5
	cmd := data & 0x1F // extract bits 4-0

	switch cmd {
	case V3_SWAP_EXACT_IN:
		ur.Command = V3SwapExactIn
	case V3_SWAP_EXACT_OUT:
		ur.Command = V3SwapExactOut
	case V2_SWAP_EXACT_IN:
		ur.Command = V2SwapExactIn
	case V2_SWAP_EXACT_OUT:
		ur.Command = V2SwapExactOut
	case PERMIT2_TRANSFER_FROM:
		ur.Command = Permit2TransferFrom
	case PERMIT2_PERMIT_BATCH:
		ur.Command = Permit2PermitBatch
	case PERMIT2_PERMIT:
		ur.Command = Permit2Permit
	case SUDOSWAP:
		ur.Command = SudoSwap
	}
	ur.CanRevert = flag == 1
	return nil
}
