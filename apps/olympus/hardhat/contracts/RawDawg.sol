// SPDX-License-Identifier: UNLICENSED

pragma solidity >=0.6.2 <0.8.20;

import "./interface/IUniswapV2Pair.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import '@uniswap/v3-periphery/contracts/libraries/TransferHelper.sol';

contract Rawdawg is Ownable {

    function executeSwap(
        address _pair,
        address _token_in,
        uint256 _amountIn,
        uint256 _amountOut,
        bool _isToken0
    ) external {
        TransferHelper.safeTransferFrom(_token_in,address(this), _pair, _amountIn);
        TransferHelper.safeApprove(_token_in, _pair, _amountIn);
        // Execute swap
        if (_isToken0) {
            IUniswapV2Pair(_pair).swap(0, _amountOut, address(this), new bytes(0));
        } else {
            IUniswapV2Pair(_pair).swap(_amountOut, 0, address(this), new bytes(0));
        }
    }
}
