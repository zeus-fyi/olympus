// SPDX-License-Identifier: UNLICENSED

pragma solidity >=0.6.2 <0.8.20;

import "./interface/IUniswapV2Pair.sol";
import "./lib/SafeTransfer.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract RawSwap is Ownable {
    using SafeTransfer for IERC20;

    function executeSwap(
        address _pair,
        address _token_in,
        uint256 _amountIn,
        uint256 _amountOut,
        bool _isToken0
    ) external {
        // Execute swap
        require(IERC20(_token_in).transfer(_pair, _amountIn));
        if (_isToken0) {
            IUniswapV2Pair(_pair).swap(0, _amountOut, address(this), new bytes(0));
        } else {
            IUniswapV2Pair(_pair).swap(_amountOut, 0, address(this), new bytes(0));
        }
    }
}
