package web3_client

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/rs/zerolog/log"
)

func NewDecodedUniversalRouterExecCmdFromMap(m map[string]interface{}) (UniversalRouterExecCmd, error) {
	log.Info().Msg("NewDecodedUniversalRouterExecCmdFromMap")
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
		subCmd.Inputs = inputsVal[i]
		err := subCmd.DecodeCommand(byteSize, inputsVal[i])
		if err != nil {
			log.Err(err).Msg("could not decode command")
			return cmds, err
		}
		subCmds[i] = subCmd
	}
	cmds.Commands = subCmds
	return cmds, nil
}

func (ur *UniversalRouterExecSubCmd) DecodeCmdByte(command byte) (bool, uint8, error) {
	buf := bytes.NewBuffer([]byte{command})
	var data uint8
	err := binary.Read(buf, binary.BigEndian, &data)
	if err != nil {
		return false, 0, err
	}
	flag := (data & 0x80) >> 7 // extract bit 7
	//ref := (data & 0x60) >> 5  // extract bits 6-5
	cmd := data & 0x1F // extract bits 4-0
	return flag == 1, cmd, nil
}

func (ur *UniversalRouterExecSubCmd) DecodeCommand(command byte, args []byte) error {
	ur.Inputs = args
	flag, cmd, err := ur.DecodeCmdByte(command)
	if err != nil {
		return err
	}
	switch cmd {
	case V3_SWAP_EXACT_IN:
		log.Info().Msg("DecodeCommand V3_SWAP_EXACT_IN")
		ur.Command = V3SwapExactIn
		params := V3SwapExactInParams{}
		err = params.Decode(ctx, ur.Inputs)
		if err != nil {
			return err
		}
		ur.DecodedInputs = params
	case V3_SWAP_EXACT_OUT:
		log.Info().Msg("DecodeCommand V3_SWAP_EXACT_OUT")
		params := V3SwapExactOutParams{}
		err = params.Decode(ctx, ur.Inputs)
		if err != nil {
			return err
		}
		ur.DecodedInputs = params
		ur.Command = V3SwapExactOut
	case V2_SWAP_EXACT_IN:
		log.Info().Msg("DecodeCommand V2_SWAP_EXACT_IN")
		params := V2SwapExactInParams{}
		err = params.Decode(ctx, ur.Inputs)
		if err != nil {
			return err
		}
		ur.DecodedInputs = params
		ur.Command = V2SwapExactIn
	case V2_SWAP_EXACT_OUT:
		log.Info().Msg("DecodeCommand V2_SWAP_EXACT_OUT")
		params := V2SwapExactOutParams{}
		err = params.Decode(ctx, ur.Inputs)
		if err != nil {
			return err
		}
		ur.DecodedInputs = params
		ur.Command = V2SwapExactOut
	case PERMIT2_TRANSFER_FROM_BATCH:
		log.Info().Msg("DecodeCommand PERMIT2_TRANSFER_FROM_BATCH")
		// TODO
		//params := Permit2PermitTransferFromBatchParams{}
		//err = params.Decode(ctx, ur.Inputs)
		//if err != nil {
		//	return err
		//}
		//ur.DecodedInputs = params
		ur.Command = Permit2TransferFromBatch
	case PERMIT2_TRANSFER_FROM:
		log.Info().Msg("DecodeCommand PERMIT2_TRANSFER_FROM")
		params := Permit2TransferFromParams{}
		err = params.Decode(ctx, ur.Inputs)
		if err != nil {
			return err
		}
		ur.DecodedInputs = params
		ur.Command = Permit2TransferFrom
	case PERMIT2_PERMIT_BATCH:
		log.Info().Msg("DecodeCommand PERMIT2_PERMIT_BATCH")
		// TODO
		//params := Permit2PermitBatchParams{}
		//err = params.Decode(ctx, ur.Inputs)
		//if err != nil {
		//	return err
		//}
		ur.Command = Permit2PermitBatch
	case PERMIT2_PERMIT:
		log.Info().Msg("DecodeCommand PERMIT2_PERMIT")
		params := Permit2PermitParams{}
		err = params.Decode(ctx, ur.Inputs)
		if err != nil {
			return err
		}
		ur.DecodedInputs = params
		ur.Command = Permit2Permit
	case SUDOSWAP:
		log.Info().Msg("DecodeCommand SUDOSWAP")
		params := SudoSwapParams{}
		err = params.Decode(ctx, ur.Inputs)
		if err != nil {
			return err
		}
		ur.DecodedInputs = params
		ur.Command = SudoSwap
	}
	ur.CanRevert = flag
	return nil
}
