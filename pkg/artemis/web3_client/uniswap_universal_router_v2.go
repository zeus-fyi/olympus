package web3_client

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
)

/*
V2_SWAP_EXACT_IN
address The recipient of the output of the trade
uint256 The amount of input tokens for the trade
uint256 The minimum amount of output tokens the user wants
address[] The UniswapV2 token path to trade along
bool A flag for whether the input tokens should come from the msg.sender (through Permit2) or whether the funds are already in the UniversalRouter

V2_SWAP_EXACT_OUT
address The recipient of the output of the trade
uint256 The amount of output tokens to receive
uint256 The maximum number of input tokens that should be spent
address[] The UniswapV2 token path to trade along
bool A flag for whether the input tokens should come from the msg.sender (through Permit2) or whether the funds are already in the UniversalRouter
*/

type V2SwapExactInParams struct {
	AmountIn        *big.Int           `json:"amountIn"`
	AmountOutMin    *big.Int           `json:"amountOutMin"`
	Path            []accounts.Address `json:"path"`
	To              accounts.Address   `json:"to"`
	InputFromSender bool               `json:"inputFromSender"`
}

type JSONV2SwapExactInParams struct {
	AmountIn        string             `json:"amountIn"`
	AmountOutMin    string             `json:"amountOutMin"`
	Path            []accounts.Address `json:"path"`
	To              accounts.Address   `json:"to"`
	InputFromSender bool               `json:"inputFromSender"`
}
