package web3_client

import (
	"fmt"
	"math/big"
)

const uniswapPriceFeeConstant = 0.3 / 100

type TradeOutcome struct {
	AmountIn   *big.Int
	AmountFees *big.Int
	AmountOut  *big.Int
}

func (p *UniswapV2Pair) PriceImpactToken1BuyToken0(tokenOneBuyAmount *big.Int) (TradeOutcome, *big.Int, *big.Int) {
	to := TradeOutcome{
		AmountIn:   tokenOneBuyAmount,
		AmountFees: nil,
		AmountOut:  nil,
	}
	tokenOneFeeDivisor := new(big.Int).Mul(tokenOneBuyAmount, big.NewInt(1000))
	tokenOneFeeDividend := new(big.Int).Mul(tokenOneBuyAmount, big.NewInt(3))
	tokenOneMinusFees := new(big.Int).Sub(tokenOneFeeDivisor, tokenOneFeeDividend)
	tokenOneFees := new(big.Int).Div(tokenOneMinusFees, big.NewInt(1000))
	feeTokenOneNormalized := new(big.Int).Sub(tokenOneBuyAmount, tokenOneFees)
	// From example: 3 Token A * 0.3% fee = 0.009 Token A
	// From example: 1200 Token A / 400 Token B = 3
	to.AmountFees = tokenOneFees

	dividend := new(big.Int).Mul(tokenOneBuyAmount, p.Reserve0)
	divisor := new(big.Int).Mul(big.NewInt(1), p.Reserve1)
	if divisor.Cmp(dividend) == 1 {
		// TODO verify this is correct
		dividend = new(big.Int).Mul(big.NewInt(1), p.Reserve0)
		divisor = new(big.Int).Mul(tokenOneBuyAmount, p.Reserve1)
	}
	tokenZeroReturnedInt := new(big.Int).Quo(dividend, divisor)
	// From example: 3 Token A * (1 Token B / 3 Token A) = 1 Token B
	fmt.Println("tokenZeroReturnedBeforeFee", tokenZeroReturnedInt.String())
	part1 := new(big.Int).Mul(tokenZeroReturnedInt, big.NewInt(3))
	part2 := new(big.Int).Mul(tokenZeroReturnedInt, big.NewInt(1000))

	fee := new(big.Int).Sub(part2, part1)
	afterFeeZero := new(big.Int).Div(fee, big.NewInt(1000))
	fmt.Println("afterFee", afterFeeZero.String())
	to.AmountOut = afterFeeZero
	feeZeroNormalized := new(big.Int).Sub(tokenZeroReturnedInt, afterFeeZero)
	fmt.Println("feeZeroNormalized", feeZeroNormalized.String())
	p.Reserve0 = new(big.Int).Sub(p.Reserve0, tokenZeroReturnedInt)
	p.Reserve0 = new(big.Int).Add(p.Reserve0, feeZeroNormalized)
	p.Reserve1 = new(big.Int).Add(p.Reserve1, tokenOneBuyAmount)
	p.Reserve1 = new(big.Int).Add(p.Reserve1, feeTokenOneNormalized)
	fmt.Println("reserve0", p.Reserve0.String())
	fmt.Println("reserve1", p.Reserve1.String())
	return to, p.Reserve0, p.Reserve1
}

func (p *UniswapV2Pair) PriceImpactToken0BuyToken1(tokenZeroBuyAmount *big.Int) (TradeOutcome, *big.Int, *big.Int) {
	to := TradeOutcome{
		AmountIn:   tokenZeroBuyAmount,
		AmountFees: nil,
		AmountOut:  nil,
	}
	tokenZeroFeeDivisor := new(big.Int).Mul(tokenZeroBuyAmount, big.NewInt(1000))
	tokenZeroFeeDividend := new(big.Int).Mul(tokenZeroBuyAmount, big.NewInt(3))
	tokenZeroMinusFees := new(big.Int).Sub(tokenZeroFeeDivisor, tokenZeroFeeDividend)
	tokenZeroFees := new(big.Int).Div(tokenZeroMinusFees, big.NewInt(1000))
	feeTokenZeroNormalized := new(big.Int).Sub(tokenZeroBuyAmount, tokenZeroFees)
	fmt.Println("feeTokenZeroNormalized", feeTokenZeroNormalized.String())
	// From example: 3 Token A * 0.3% fee = 0.009 Token A
	// From example: 1200 Token A / 400 Token B = 3
	to.AmountFees = tokenZeroFees
	fmt.Println("tokenZeroBuyAmount", tokenZeroBuyAmount.String())
	fmt.Println("reserve0", p.Reserve0.String())
	fmt.Println("reserve1", p.Reserve1.String())

	dividend := new(big.Int).Mul(tokenZeroBuyAmount, p.Reserve1)
	divisor := new(big.Int).Mul(big.NewInt(1), p.Reserve0)
	if divisor.Cmp(dividend) == 1 {
		// TODO verify this is correct
		dividend = new(big.Int).Mul(big.NewInt(1), p.Reserve1)
		divisor = new(big.Int).Mul(tokenZeroBuyAmount, p.Reserve0)
	}
	fmt.Println("dividend", dividend.String())
	fmt.Println("divisor", divisor.String())
	tokenOneReturnedInt := new(big.Int).Quo(dividend, divisor)
	fmt.Println("tokenOneReturnedInt", tokenOneReturnedInt.String())
	// From example: 3 Token A * (1 Token B / 3 Token A) = 1 Token B
	fmt.Println("tokenOneReturnedBeforeFee", tokenOneReturnedInt.String())
	part1 := new(big.Int).Mul(tokenOneReturnedInt, big.NewInt(3))
	part2 := new(big.Int).Mul(tokenOneReturnedInt, big.NewInt(1000))

	fee := new(big.Int).Sub(part2, part1)
	afterFeeOne := new(big.Int).Div(fee, big.NewInt(1000))
	fmt.Println("afterFee", afterFeeOne.String())
	to.AmountOut = afterFeeOne
	feeOneNormalized := new(big.Int).Sub(tokenOneReturnedInt, afterFeeOne)
	fmt.Println("feeOneNormalized", feeOneNormalized.String())
	p.Reserve1 = new(big.Int).Sub(p.Reserve1, tokenOneReturnedInt)
	p.Reserve1 = new(big.Int).Add(p.Reserve1, feeOneNormalized)
	p.Reserve0 = new(big.Int).Add(p.Reserve0, tokenZeroBuyAmount)
	p.Reserve0 = new(big.Int).Add(p.Reserve0, feeTokenZeroNormalized)
	return to, p.Reserve0, p.Reserve1
}
