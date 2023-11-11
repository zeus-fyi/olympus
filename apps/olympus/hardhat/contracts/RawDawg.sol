// SPDX-License-Identifier: UNLICENSED

pragma solidity >=0.6.2 <0.8.20;

import "./interface/IUniswapV2Pair.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import '@uniswap/v3-periphery/contracts/libraries/TransferHelper.sol';
import '@uniswap/universal-router/contracts/interfaces/IUniversalRouter.sol';

contract Rawdawg is Ownable {
    address public constant universalRouterAddress = 0xEf1c6E67703c7BD7107eed8303Fbe6EC2554BF6B;
    address public constant v2routerAddress = 0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D;

    receive() external payable {}
    fallback() external payable {}

    function executeUniversalRouter(
        bytes calldata commands,
        bytes[] calldata inputs,
        uint256 deadline
    ) external payable {
        if (msg.value > 0) {
            (bool sent, ) = payable(universalRouterAddress).call{value: msg.value}("");
            require(sent, "Failed to send Ether");
        }
        IUniversalRouter(universalRouterAddress).execute(commands, inputs, deadline);
    }

    struct swapParams {
        address _pair;
        address _token_in;
        bool _isToken0;
        uint256 _amountIn;
        uint256 _amountOut;
    }

//    function batchExecuteSwap(
//        swapParams[] calldata _swap
//    ) external {
//        uint256 length = _swap.length;
//        for (uint256 i = 0; i < length;) {
//            _executeSwap(_swap[i]._pair, _swap[i]._token_in, _swap[i]._amountIn, _swap[i]._amountOut, _swap[i]._isToken0);
//            unchecked {
//                ++i;
//            }
//        }
//    }

    function executeSwap(
        address _pair,
        address _token_in,
        bool _isToken0,
        uint256 _amountIn,
        uint256 _amountOut
    ) external {
        _executeSwap(_pair, _token_in, _isToken0, _amountIn, _amountOut);
    }

    function _executeSwap(
        address _pair,
        address _token_in,
        bool _isToken0,
        uint256 _amountIn,
        uint256 _amountOut
    ) internal {
        TransferHelper.safeTransfer(_token_in, _pair, _amountIn);
        // wondering if just two functions eg. swapTokenZero, or swapTokenOne is better?
        TransferHelper.safeApprove(_token_in, v2routerAddress, _amountIn);
        // Execute swap
        if (_isToken0) {
            IUniswapV2Pair(_pair).swap(0, _amountOut, address(this), new bytes(0));
        } else {
            IUniswapV2Pair(_pair).swap(_amountOut, 0, address(this), new bytes(0));
        }
    }

    // Function to simulate a swap and then revert
    function simulateV2AndRevertSwap(
        address _pair,
        address _token_in,
        address _token_out,
        bool _isToken0,
        uint256 _amountIn,
        uint256 _amountOut
    ) external returns (uint256 balanceTokenInAfter, uint256 balanceTokenOutAfter, uint256 gasUsed) {
        uint256 gasBefore = gasleft();

        try
            this.executeSwap(_pair, _token_in, _isToken0, _amountIn, _amountOut) {
            // If the swap is successful, we immediately revert
            revert("Simulation only - reverting transaction");
        } catch {
            // Calculate the gas used
            gasUsed = gasBefore - gasleft();

            // Calculate the simulated final balances
            balanceTokenInAfter = IERC20(_token_in).balanceOf(address(this));
            balanceTokenOutAfter = IERC20(_token_out).balanceOf(address(this));
        }
    }
}


//function handleRevert(
//        bytes memory reason,
//        IUniswapV3Pool pool,
//        uint256 gasEstimate
//    )
//        private
//        view
//        returns (
//            uint256 amount,
//            uint160 sqrtPriceX96After,
//            uint32 initializedTicksCrossed,
//            uint256
//        )
//    {
//        int24 tickBefore;
//        int24 tickAfter;
//        (, tickBefore, , , , , ) = pool.slot0();
//        (amount, sqrtPriceX96After, tickAfter) = parseRevertReason(reason);
//
//        initializedTicksCrossed = pool.countInitializedTicksCrossed(tickBefore, tickAfter);
//
//        return (amount, sqrtPriceX96After, initializedTicksCrossed, gasEstimate);
//    }

//    function parseRevertReason(bytes memory reason)
//    private
//    pure
//    returns (
//        uint256 amount,
//        uint160 sqrtPriceX96After,
//        int24 tickAfter
//    )
//    {
//        if (reason.length != 96) {
//            if (reason.length < 68) revert('Unexpected error');
//            assembly {
//                reason := add(reason, 0x04)
//            }
//            revert(abi.decode(reason, (string)));
//        }
//        return abi.decode(reason, (uint256, uint160, int24));
//    }