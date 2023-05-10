package web3_client

import (
	"context"
	"errors"
	"math/big"

	"github.com/gochain/gochain/v4/common"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

type UniswapV2Pair struct {
	PairContractAddr     string
	Price0CumulativeLast *big.Int
	Price1CumulativeLast *big.Int
	KLast                *big.Int
	Token0               common.Address
	Token1               common.Address
	Reserve0             *big.Int
	Reserve1             *big.Int
	BlockTimestampLast   *big.Int
}

func (u *UniswapV2Client) GetPairContractPrices(ctx context.Context, pairContractAddr string) (UniswapV2Pair, error) {
	//addrOne, addrTwo := StringsToAddresses(addressOne, addressTwo)
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: pairContractAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       u.PairAbi,
	}
	pairInfo := UniswapV2Pair{
		PairContractAddr:     pairContractAddr,
		Price0CumulativeLast: nil,
		Price1CumulativeLast: nil,
		KLast:                nil,
		Token0:               common.Address{},
		Token1:               common.Address{},
		Reserve0:             nil,
		Reserve1:             nil,
		BlockTimestampLast:   nil,
	}
	price0, err := u.SingleReadMethodBigInt(ctx, "price0CumulativeLast", scInfo)
	if err != nil {
		return UniswapV2Pair{}, err
	}
	pairInfo.Price0CumulativeLast = price0
	price1, err := u.SingleReadMethodBigInt(ctx, "price1CumulativeLast", scInfo)
	if err != nil {
		return UniswapV2Pair{}, err
	}
	pairInfo.Price1CumulativeLast = price1
	kLast, err := u.SingleReadMethodBigInt(ctx, "kLast", scInfo)
	if err != nil {
		return UniswapV2Pair{}, err
	}
	pairInfo.KLast = kLast
	token0, err := u.SingleReadMethodAddr(ctx, "token0", scInfo)
	if err != nil {
		return UniswapV2Pair{}, err
	}
	pairInfo.Token0 = token0
	token1, err := u.SingleReadMethodAddr(ctx, "token1", scInfo)
	if err != nil {
		return UniswapV2Pair{}, err
	}
	pairInfo.Token1 = token1
	scInfo.MethodName = "getReserves"
	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return UniswapV2Pair{}, err
	}
	if len(resp) <= 2 {
		return UniswapV2Pair{}, err
	}
	reserve0, err := ParseBigInt(resp[0])
	if err != nil {
		return UniswapV2Pair{}, err
	}
	pairInfo.Reserve0 = reserve0
	reserve1, err := ParseBigInt(resp[1])
	if err != nil {
		return UniswapV2Pair{}, err
	}
	pairInfo.Reserve1 = reserve1
	blockTimestampLast, err := ParseBigInt(resp[2])
	if err != nil {
		return UniswapV2Pair{}, err
	}
	pairInfo.BlockTimestampLast = blockTimestampLast
	return pairInfo, nil
}

func (u *UniswapV2Client) SingleReadMethodBigInt(ctx context.Context, methodName string, scInfo *web3_actions.SendContractTxPayload) (*big.Int, error) {
	scInfo.MethodName = methodName
	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return &big.Int{}, err
	}
	if len(resp) == 0 {
		return &big.Int{}, errors.New("empty response")
	}
	bi, err := ParseBigInt(resp[0])
	if err != nil {
		return &big.Int{}, err
	}
	return bi, nil
}

func (u *UniswapV2Client) SingleReadMethodAddr(ctx context.Context, methodName string, scInfo *web3_actions.SendContractTxPayload) (common.Address, error) {
	scInfo.MethodName = methodName
	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return common.Address{}, err
	}
	if len(resp) == 0 {
		return common.Address{}, errors.New("empty response")
	}
	addr, err := ConvertToAddress(resp[0])
	if err != nil {
		return common.Address{}, err
	}
	return addr, nil
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
