// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "./interface/IERC20.sol";

contract Rawdawg {
    address internal immutable user;

    // transfer(address,uint256)
    bytes4 internal constant ERC20_TRANSFER_ID = 0xa9059cbb;
    // transferFrom(address,address,uint256)
    bytes4 internal constant ERC20_TRANSFER_FROM = 0x23b872dd;
    // SafeTransferFrom(address,address,uint256)
    bytes4 internal constant ERC721_SAFE_TRANSFER_FROM = 0x42842e0e;
    // Uniswapv2 Pair Swap(unint256,uint256,address,bytes)
    bytes4 internal constant PAIR_SWAP_ID = 0x022c0d9f;

    receive() external payable {}
    constructor() {
        user = msg.sender;
    }

    function drainERC20(address token) public {
        require(msg.sender == user, "RawDawg: Only the user can drain ERC20s");
        bytes memory payload = abi.encodeWithSelector(ERC20_TRANSFER_ID, user, IERC20(token).balanceOf(address(this)));
        (bool success, ) = token.call(payload);
        require(success, "RawDawg: ERC20 transfer failed");
    }

    fallback() external payable {
        require(msg.sender == user, "RawDawg: Only the user can call fallback");
        address memUser = user;

        assembly {
            if iszero(eq(caller(), memUser)){
                revert(3, 3)
            }
            // bytes20
            let token := shr(96, calldataload(0x00))
            // bytes20
            let pair  := shr(96, calldataload(0x14))
            // unit128
            let amountIn := shr(128, calldataload(0x28))
            // unit128
            let amountOut := shr(128, calldataload(0x38))
            // unit8
            let tokenOutNo := shr(248, calldataload(0x48))
            // Todo: maybe add another variable for uiniswapv2 function selector

            // call token.transfer(pair, amountIn)
            mstore(0x7c, ERC20_TRANSFER_ID)
            mstore(0x80, pair)
            mstore(0xa0, amountIn)
            let success1 := call(gas(), token, 0, 0x7c, 0xc0, 0x0, 0x0)
            if iszero(success1) {
                revert(3, 3)
            }

            /* call pair.swap(
                tokenOutNo == 0 ? amountOut : 0,
                tokenOutNo == 1 ? amountOut : 0,
                address(this),
                new bytes(0)
            )*/
            mstore(0x7c, PAIR_SWAP_ID)
            switch tokenOutNo
            case 0 {
                mstore(0x80, amountOut)
                mstore(0xa0, 0)
            }
            case 1 {
                mstore(0x80, 0)
                mstore(0xa0, amountOut)
            }
            // address(this)
            mstore(0xc0, address())
            // empty bytes
            mstore(0xe0, 0x80)
            let success2 := call(gas(), pair, 0, 0x7c, 0xa4, 0, 0)
            if iszero(success2) {
                revert(3, 3)
            }
        }
    }
}