package web3_client

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
)

func (ur *UniversalRouterExecCmd) EncodeCommands(ctx context.Context, abiFile *abi.ABI) (*UniversalRouterExecParams, error) {
	//log.Info().Msg("UniversalRouterExecCmd: EncodeCommands")
	encodedCmd := &UniversalRouterExecParams{}
	for _, cmd := range ur.Commands {
		cmdByteStr, inputs, err := cmd.EncodeCommand(ctx, abiFile)
		if err != nil {
			log.Err(err).Msg("could not encode command")
			return nil, err
		}
		encodedCmd.Commands = append(encodedCmd.Commands, cmdByteStr)
		encodedCmd.Inputs = append(encodedCmd.Inputs, inputs)
	}
	encodedCmd.Deadline = ur.Deadline
	encodedCmd.Payable = ur.Payable
	return encodedCmd, nil
}

func (ur *UniversalRouterExecSubCmd) EncodeCommandByte(flag bool, command int) byte {
	data := command // start with command in the lower 5 bits
	if flag {       // if flag is true, set the highest bit
		data |= 0x80
	}
	return byte(data)
}

func (ur *UniversalRouterExecSubCmd) EncodeCommand(ctx context.Context, abiFile *abi.ABI) (byte, []byte, error) {
	var cmdByte byte
	switch ur.Command {
	case V3SwapExactIn:
		//log.Info().Msg("EncodeCommand V3_SWAP_EXACT_IN")
		params := ur.DecodedInputs.(V3SwapExactInParams)
		inputs, err := params.Encode(ctx, abiFile)
		if err != nil {
			return cmdByte, nil, err
		}
		cmdByte = ur.EncodeCommandByte(ur.CanRevert, V3_SWAP_EXACT_IN)
		return cmdByte, inputs, nil
	case V3SwapExactOut:
		//log.Info().Msg("EncodeCommand V3_SWAP_EXACT_OUT")
		params := ur.DecodedInputs.(V3SwapExactOutParams)
		inputs, err := params.Encode(ctx, abiFile)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		cmdByte = ur.EncodeCommandByte(ur.CanRevert, V3_SWAP_EXACT_OUT)
		return cmdByte, inputs, nil
	case V2SwapExactIn:
		//log.Info().Msg("EncodeCommand V2_SWAP_EXACT_IN")
		params := ur.DecodedInputs.(V2SwapExactInParams)
		inputs, err := params.Encode(ctx, abiFile)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		cmdByte = ur.EncodeCommandByte(ur.CanRevert, V2_SWAP_EXACT_IN)
		return cmdByte, inputs, nil
	case V2SwapExactOut:
		//log.Info().Msg("EncodeCommand V2_SWAP_EXACT_OUT")
		params := ur.DecodedInputs.(V2SwapExactOutParams)
		inputs, err := params.Encode(ctx, abiFile)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		cmdByte = ur.EncodeCommandByte(ur.CanRevert, V2_SWAP_EXACT_OUT)
		return cmdByte, inputs, nil
	case Permit2TransferFromBatch:
		//	log.Info().Msg("EncodeCommand PERMIT2_TRANSFER_FROM_BATCH")
		params := ur.DecodedInputs.(Permit2PermitTransferFromBatchParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		cmdByte = ur.EncodeCommandByte(ur.CanRevert, PERMIT2_TRANSFER_FROM_BATCH)
		ur.Command = Permit2TransferFromBatch
		return cmdByte, inputs, nil
	case Transfer:
		params := ur.DecodedInputs.(TransferParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		cmdByte = ur.EncodeCommandByte(ur.CanRevert, TRANSFER)
		return cmdByte, nil, nil
	case Permit2TransferFrom:
		//log.Info().Msg("EncodeCommand PERMIT2_TRANSFER_FROM")
		params := ur.DecodedInputs.(Permit2TransferFromParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		cmdByte = ur.EncodeCommandByte(ur.CanRevert, PERMIT2_TRANSFER_FROM)
		return cmdByte, inputs, nil
	case Permit2PermitBatch:
		//log.Info().Msg("EncodeCommand PERMIT2_PERMIT_BATCH")
		params := ur.DecodedInputs.(Permit2PermitBatchParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		cmdByte = ur.EncodeCommandByte(ur.CanRevert, PERMIT2_PERMIT_BATCH)
		return cmdByte, inputs, nil
	case Permit2Permit:
		//log.Info().Msg("EncodeCommand PERMIT2_PERMIT")
		params := ur.DecodedInputs.(Permit2PermitParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		cmdByte = ur.EncodeCommandByte(ur.CanRevert, PERMIT2_PERMIT)
		return cmdByte, inputs, nil
	case SudoSwap:
		//log.Info().Msg("EncodeCommand SUDOSWAP")
		params := ur.DecodedInputs.(SudoSwapParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, inputs, err
		}
		ur.Inputs = inputs
		cmdByte = ur.EncodeCommandByte(ur.CanRevert, SUDOSWAP)
		return cmdByte, inputs, nil
	case WrapETH:
		//log.Info().Msg("EncodeCommand WRAP_ETH")
		params := ur.DecodedInputs.(WrapETHParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, inputs, err
		}
		ur.Inputs = inputs
		cmdByte = ur.EncodeCommandByte(ur.CanRevert, WRAP_ETH)
		return cmdByte, inputs, nil
	case UnwrapWETH:
		//log.Info().Msg("EncodeCommand UNWRAP_WETH")
		params := ur.DecodedInputs.(UnwrapWETHParams)
		inputs, err := params.Encode(ctx, abiFile)
		if err != nil {
			return cmdByte, inputs, err
		}
		ur.Inputs = inputs
		cmdByte = ur.EncodeCommandByte(ur.CanRevert, UNWRAP_WETH)
		return cmdByte, inputs, nil
	default:
	}
	return cmdByte, nil, errors.New("unknown command")
}
