// SPDX-License-Identifier: UNLICENSED

pragma solidity >=0.6.2 <0.8.20;

import "./interface/IUniswapV2Pair.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import '@uniswap/v3-periphery/contracts/libraries/TransferHelper.sol';
import '@uniswap/v3-periphery/contracts/interfaces/IQuoterV2.sol';
import '@uniswap/universal-router/contracts/interfaces/IUniversalRouter.sol';
import "@uniswap/v3-core/contracts/interfaces/IUniswapV3Pool.sol";
import "@uniswap/v2-periphery/contracts/interfaces/IUniswapV2Router02.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract Rawdawg is Ownable {
    address public constant universalRouterAddress = 0xEf1c6E67703c7BD7107eed8303Fbe6EC2554BF6B;
    address public constant routerV2Address = 0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D;
    address public constant quoterV2Address = 0x61fFE014bA17989E743c5F6cB21bF9697530B21e;
    uint160 internal constant MIN_SQRT_RATIO = 4295128739;
    uint160 internal constant MAX_SQRT_RATIO = 1461446703485210103287273052203988822378723970342;

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
        TransferHelper.safeApprove(_token_in, routerV2Address, _amountIn);
        // Execute swap
        if (_isToken0) {
            IUniswapV2Pair(_pair).swap(0, _amountOut, address(this), new bytes(0));
        } else {
            IUniswapV2Pair(_pair).swap(_amountOut, 0, address(this), new bytes(0));
        }
    }

    function simulateV2AndRevertSwap(
        address _pair,
        address _token_in,
        address _token_out,
        bool _isToken0,
        uint256 _amountIn,
        uint256 _amountOut
    ) external {
        IUniswapV2Router02 router = IUniswapV2Router02(routerV2Address);

        uint256 buyAmountOut;
        uint256 buyAmountOutExpected;
        uint256 buyGas;

    {
            uint256 _amountIn = IERC20(_token_in).balanceOf(address(this));
            require(_amountIn > 0, "Checker: BUY_INPUT_ZERO");
            address[] memory pathBuy  = new address[](2);
            pathBuy[0] = _token_in;
            pathBuy[1] = _token_out;
            buyAmountOutExpected = router.getAmountsOut(_amountIn, pathBuy)[1];
            uint256 startBuyGas = gasleft();
            try router.swapExactTokensForTokensSupportingFeeOnTransferTokens(_amountIn, 0, pathBuy, address(this), block.timestamp+500){
                buyGas = startBuyGas - gasleft();
                buyAmountOut = IERC20(_token_out).balanceOf(address(this)); // - tokensBefore;
                require(buyAmountOut > 0, "Checker: BUY_OUTPUT_ZERO");
            }
            catch Error(string memory){
                revert("Checker: BUY_FAILED");
            }
        }
        /*
           return [
               buyAmountOut,
               buyAmountOutExpected,
               sellAmountOut,
               sellAmountOutExpected,
               buyGas,
               sellGas
           ];
       */

        // returns(uint256[3] memory)
        assembly{
            let ret := mload(0x40)
            mstore(ret,            buyAmountOut)
            mstore(add(ret, 0x20), buyAmountOutExpected)
            mstore(add(ret, 0x40), buyGas)
            return(ret, mul(0x20, 3))
        }
    }
}

