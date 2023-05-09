package web3_client

import "github.com/gochain/gochain/v4/common"

func (u *UniswapV2Client) GetPairContractPrices(addressOne, addressTwo common.Address) {
	// TODO, smart contract get pair call
	// price0CumulativeLast
	// price1CumulativeLast
	// token0
	// token1
	// getReserves
	return
}

// TODO
// function swap(uint amount0Out, uint amount1Out, address to, bytes calldata data) external;

// get k value, x*y=k.
//   uint public kLast; // reserve0 * reserve1, as of immediately after the most recent liquidity event
