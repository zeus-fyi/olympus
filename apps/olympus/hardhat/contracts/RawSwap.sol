// SPDX-License-Identifier: UNLICENSED

pragma solidity >=0.6.2 <0.8.20;

import "./interface/IUniswapV2Pair.sol";
import "./lib/SafeTransfer.sol";

contract RawSwap {
    using SafeTransfer for IERC20;

    function executeSwap(
        address _pair,
        uint256 _amount0Out,
        uint256 _amount1Out
    ) external {
        // Execute swap
        IUniswapV2Pair(_pair).swap(_amount0Out, _amount1Out, address(this), new bytes(0));
    }
}

