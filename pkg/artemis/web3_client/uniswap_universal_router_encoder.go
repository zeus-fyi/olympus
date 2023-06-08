package web3_client

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
)

func (ur *UniversalRouterExecCmd) EncodeCommands(ctx context.Context) (*UniversalRouterExecParams, error) {
	log.Info().Msg("UniversalRouterExecCmd: EncodeCommands")
	encodedCmd := &UniversalRouterExecParams{}
	for _, cmd := range ur.Commands {
		cmdByteStr, inputs, err := cmd.EncodeCommand(ctx)
		if err != nil {
			log.Err(err).Msg("could not encode command")
			return nil, err
		}
		encodedCmd.Commands = append(encodedCmd.Commands, cmdByteStr)
		encodedCmd.Inputs = append(encodedCmd.Inputs, inputs)
	}
	return encodedCmd, nil
}

func (ur *UniversalRouterExecSubCmd) EncodeCommandByte() byte {
	var cmdByte byte
	return cmdByte
}

func (ur *UniversalRouterExecSubCmd) EncodeCommand(ctx context.Context) (byte, []byte, error) {
	var cmdByte byte
	switch ur.Command {
	case V3SwapExactIn:
		log.Info().Msg("EncodeCommand V3_SWAP_EXACT_IN")
		params := ur.DecodedInputs.(V3SwapExactInParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, nil, err
		}
		return cmdByte, inputs, nil
	case V3SwapExactOut:
		log.Info().Msg("DecodeCommand V3_SWAP_EXACT_OUT")
		params := ur.DecodedInputs.(V3SwapExactOutParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		return cmdByte, inputs, nil
	case V2SwapExactIn:
		log.Info().Msg("DecodeCommand V2_SWAP_EXACT_IN")
		params := ur.DecodedInputs.(V2SwapExactInParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		return cmdByte, inputs, nil
	case V2SwapExactOut:
		log.Info().Msg("DecodeCommand V2_SWAP_EXACT_OUT")
		params := ur.DecodedInputs.(V2SwapExactOutParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		return cmdByte, inputs, nil
	case Permit2TransferFromBatch:
		log.Info().Msg("DecodeCommand PERMIT2_TRANSFER_FROM_BATCH")
		// TODO
		//params := ur.DecodedInputs.(Permit2PermitTransferFromBatchParams)
		//err = params.Decode(ctx, ur.Inputs)
		//if err != nil {
		//	return err
		//}
		//ur.DecodedInputs = params
		ur.Command = Permit2TransferFromBatch
	case Permit2TransferFrom:
		log.Info().Msg("DecodeCommand PERMIT2_TRANSFER_FROM")
		params := ur.DecodedInputs.(Permit2TransferFromParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		return cmdByte, inputs, nil
	case Permit2PermitBatch:
		log.Info().Msg("DecodeCommand PERMIT2_PERMIT_BATCH")
		// TODO
		//params := ur.DecodedInputs.(Permit2PermitBatchParams)

		//err = params.Decode(ctx, ur.Inputs)
		//if err != nil {
		//	return err
		//}
		ur.Command = Permit2PermitBatch
	case Permit2Permit:
		log.Info().Msg("DecodeCommand PERMIT2_PERMIT")
		params := ur.DecodedInputs.(Permit2PermitParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		return cmdByte, inputs, nil
	case SudoSwap:
		log.Info().Msg("DecodeCommand SUDOSWAP")
		params := ur.DecodedInputs.(SudoSwapParams)
		inputs, err := params.Encode(ctx)
		if err != nil {
			return cmdByte, nil, err
		}
		ur.Inputs = inputs
		return cmdByte, inputs, nil
	default:
	}
	return cmdByte, nil, errors.New("unknown command")
}
