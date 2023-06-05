package web3_client

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
)

/*
The inputs for V3_SWAP_EXACT_IN is the encoding of 5 parameters:

address The recipient of the output of the trade
uint256 The amount of input tokens for the trade
uint256 The minimum amount of output tokens the user wants
bytes The UniswapV3 path you want to trade along
bool A flag for whether the input funds should come from the caller (through Permit2) or whether the funds are already in the UniversalRouter
*/

type V3SwapExactInParams struct {
	AmountIn        *big.Int         `json:"amountIn"`
	AmountOutMin    *big.Int         `json:"amountOutMin"`
	Path            []byte           `json:"path"`
	To              accounts.Address `json:"to"`
	InputFromSender bool             `json:"inputFromSender"`
}

type JSONV3SwapExactInParams struct {
	AmountIn        string           `json:"amountIn"`
	AmountOutMin    string           `json:"amountOutMin"`
	Path            []byte           `json:"path"`
	To              accounts.Address `json:"to"`
	InputFromSender bool             `json:"inputFromSender"`
}

/*
V3_SWAP_EXACT_OUT
address The recipient of the output of the trade
uint256 The amount of output tokens to receive
uint256 The maximum number of input tokens that should be spent
bytes The UniswapV3 encoded path to trade along
bool A flag for whether the input tokens should come from the msg.sender (through Permit2) or whether the funds are already in the UniversalRouter
*/
