// SPDX-License-Identifier: UNLICENSED

pragma solidity >=0.6.2 <0.8.20;

import "./interface/IUniswapV2Pair.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import '@uniswap/v3-periphery/contracts/libraries/TransferHelper.sol';

contract Rawdawg is Ownable {
    address public constant routerAddress = 0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D;

    struct swapParams {
        address _pair;
        address _token_in;
        uint256 _amountIn;
        uint256 _amountOut;
        bool _isToken0;
    }

    function batchExecuteSwap(
        swapParams[] calldata _swap
    ) external {
        for (uint256 i = 0; i < _swap.length;) {
            _executeSwap(_swap[i]._pair, _swap[i]._token_in, _swap[i]._amountIn, _swap[i]._amountOut, _swap[i]._isToken0);
            unchecked {
                ++i;
            }
        }
    }

    function executeSwap(
        address _pair,
        address _token_in,
        uint256 _amountIn,
        uint256 _amountOut,
        bool _isToken0
    ) external {
        _executeSwap(_pair, _token_in, _amountIn, _amountOut, _isToken0);
    }

    function _executeSwap(
        address _pair,
        address _token_in,
        uint256 _amountIn,
        uint256 _amountOut,
        bool _isToken0
    ) internal {
        TransferHelper.safeTransfer(_token_in, _pair, _amountIn);
        TransferHelper.safeApprove(_token_in, routerAddress, _amountIn);
        // Execute swap
        if (_isToken0) {
            IUniswapV2Pair(_pair).swap(0, _amountOut, address(this), new bytes(0));
        } else {
            IUniswapV2Pair(_pair).swap(_amountOut, 0, address(this), new bytes(0));
        }
    }
}
