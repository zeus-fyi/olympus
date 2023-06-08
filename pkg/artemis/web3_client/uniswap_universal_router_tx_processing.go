package web3_client

import (
	"context"
)

// todo add counter

func (u *UniswapClient) ProcessUniversalRouterTxs(ctx context.Context, tx MevTx) {
	subcmd, err := NewDecodedUniversalRouterExecCmdFromMap(tx.Args)
	if err != nil {
		return
	}

	// todo, update this from stub to real
	pair := UniswapV2Pair{}

	// todo needs to save trade analysis results
	for _, subtx := range subcmd.Commands {
		switch subtx.Command {
		case V3SwapExactIn:
			inputs := subtx.DecodedInputs.(V3SwapExactInParams)
			tf := inputs.BinarySearch(pair)
			tf.InitialPair = pair.ConvertToJSONType()
		case V3SwapExactOut:
			inputs := subtx.DecodedInputs.(V3SwapExactOutParams)
			tf := inputs.BinarySearch(pair)
			tf.InitialPair = pair.ConvertToJSONType()
		case V2SwapExactIn:
			inputs := subtx.DecodedInputs.(V2SwapExactInParams)
			tf := inputs.BinarySearch(pair)
			tf.InitialPair = pair.ConvertToJSONType()
		case V2SwapExactOut:
			inputs := subtx.DecodedInputs.(V2SwapExactOutParams)
			tf := inputs.BinarySearch(pair)
			tf.InitialPair = pair.ConvertToJSONType()
		default:
		}
	}
}
