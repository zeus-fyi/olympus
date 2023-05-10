package web3_client

import (
	"fmt"
	"math/big"
)

const uniswapPriceFeeConstant = 0.3 / 100

func (p *UniswapV2Pair) PriceImpact(tokenOneBuyAmount *big.Int) *big.Float {
	tokenOneAmountFloat := new(big.Float).SetInt(tokenOneBuyAmount)
	// Calculate fee
	feeTokenOne := new(big.Float).Mul(tokenOneAmountFloat, big.NewFloat(uniswapPriceFeeConstant))
	fmt.Println("feeOne", feeTokenOne.String())
	price, err := p.GetToken1Price()
	if err != nil {
		return nil
	}
	returnBeforeFeeTokenZero := new(big.Float).Mul(tokenOneAmountFloat, price)
	fmt.Println("returnBeforeFeeTokenZero", returnBeforeFeeTokenZero.String())
	feeTokenZero := new(big.Float).Mul(returnBeforeFeeTokenZero, big.NewFloat(uniswapPriceFeeConstant))
	fmt.Println("feeTokenZero", feeTokenZero.String())
	// Update reserves
	tokenOneFeeInt, _ := feeTokenOne.Int(nil)
	p.Reserve1.Add(p.Reserve1, tokenOneFeeInt)
	p.Reserve1.Add(p.Reserve1, tokenOneBuyAmount)

	feeTokenZeroInt, _ := feeTokenZero.Int(nil)
	p.Reserve0.Add(p.Reserve0, feeTokenZeroInt)
	tokenZeroPurchaseAmount, _ := returnBeforeFeeTokenZero.Int(nil)
	p.Reserve0.Sub(p.Reserve0, tokenZeroPurchaseAmount)

	fmt.Println("reserve0", p.Reserve0.String())
	fmt.Println("reserve1", p.Reserve1.String())
	// Calculate new price
	newPrice, _ := p.GetToken1Price()
	return newPrice
}
