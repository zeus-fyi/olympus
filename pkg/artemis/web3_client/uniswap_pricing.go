package web3_client

import (
	"fmt"
	"math/big"
)

const uniswapPriceFeeConstant = 0.3 / 100

func (p *UniswapV2Pair) PriceImpact(tokenOneBuyAmount *big.Int) (*big.Float, *big.Float) {
	// From example: 3 Token A
	fmt.Println("tokenOneBuyAmount", tokenOneBuyAmount.String())
	tokenOneAmountFloat := new(big.Float).SetInt(tokenOneBuyAmount)
	// From example: 3 Token A * 0.3% fee
	feeTokenOne := new(big.Float).Mul(tokenOneAmountFloat, big.NewFloat(uniswapPriceFeeConstant))
	// From example: 3 Token A * 0.3% fee = 0.009 Token A
	fmt.Println("feeOne", feeTokenOne.String())
	// From example: 1200 Token A / 400 Token B = 3
	priceToken0, err := p.GetToken0Price()
	if err != nil {
		return nil, nil
	}
	fmt.Println("price token A per token B", priceToken0.String())
	priceToken1, err := p.GetToken1Price()
	if err != nil {
		return nil, nil
	}
	fmt.Println("price token B per token A", priceToken1.String())
	tokenZeroReturned := new(big.Float).Mul(tokenOneAmountFloat, priceToken1)

	// From example: 3 Token A * (1 Token B / 3 Token A) = 1 Token B
	fmt.Println("tokenZeroReturnedBeforeFee", tokenZeroReturned.String())
	feeTokenZero := new(big.Float).Mul(tokenZeroReturned, big.NewFloat(uniswapPriceFeeConstant))
	// From example: 1 Token B * 0.3% fee = 0.003 Token B
	fmt.Println("feeTokenZero", feeTokenZero.String())
	// Update reserves
	tokenOneFeeInt, _ := feeTokenOne.Int(nil)
	p.Reserve1.Add(p.Reserve1, tokenOneFeeInt)
	p.Reserve1.Add(p.Reserve1, tokenOneBuyAmount)

	feeTokenZeroInt, _ := feeTokenZero.Int(nil)
	p.Reserve0.Add(p.Reserve0, feeTokenZeroInt)
	tokenZeroPurchaseAmount, _ := tokenZeroReturned.Int(nil)
	p.Reserve0.Sub(p.Reserve0, tokenZeroPurchaseAmount)

	// From example: 1200 Token A + 3 Token A + 0.009 Token A = 1203.009 Token A
	fmt.Println("reserve0", p.Reserve0.String())
	// From example: 400 Token B - 1 Token B + 0.003 Token B = 399.003 Token B
	fmt.Println("reserve1", p.Reserve1.String())
	// Calculate new price
	newPriceToken1, _ := p.GetToken1Price()
	newPriceToken0, _ := p.GetToken0Price()
	return newPriceToken1, newPriceToken0
}
