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

/*
  // if fee is on, mint liquidity equivalent to 1/6th of the growth in sqrt(k)
    function _mintFee(uint112 _reserve0, uint112 _reserve1) private returns (bool feeOn) {
        address feeTo = IUniswapV2Factory(factory).feeTo();
        feeOn = feeTo != address(0);
        uint _kLast = kLast; // gas savings
        if (feeOn) {
            if (_kLast != 0) {
                uint rootK = Math.sqrt(uint(_reserve0).mul(_reserve1));
                uint rootKLast = Math.sqrt(_kLast);
                if (rootK > rootKLast) {
                    uint numerator = totalSupply.mul(rootK.sub(rootKLast));
                    uint denominator = rootK.mul(5).add(rootKLast);
                    uint liquidity = numerator / denominator;
                    if (liquidity > 0) _mint(feeTo, liquidity);
                }
            }
        } else if (_kLast != 0) {
            kLast = 0;
        }
    }
*/
// https://github.com/Uniswap/v2-core/blob/ee547b17853e71ed4e0101ccfd52e70d5acded58/contracts/UniswapV2Pair.sol#L26
// get k value, x*y=k.
//   uint public kLast; // reserve0 * reserve1, as of immediately after the most recent liquidity event
