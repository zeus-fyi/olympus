package web3_client

import (
	"errors"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/v4/common"
)

func (p *UniswapV2Pair) PriceImpact(tokenAddrPath common.Address, tokenBuyAmount *big.Int) (TradeOutcome, error) {
	tokenNumber := p.GetTokenNumber(tokenAddrPath)
	switch tokenNumber {
	case 1:
		to, _, _ := p.PriceImpactToken1BuyToken0(tokenBuyAmount)
		to.AmountInAddr = tokenAddrPath
		to.AmountOutAddr = p.GetOppositeToken(tokenAddrPath.String())
		return to, nil
	case 0:
		to, _, _ := p.PriceImpactToken0BuyToken1(tokenBuyAmount)
		to.AmountInAddr = tokenAddrPath
		to.AmountOutAddr = p.GetOppositeToken(tokenAddrPath.String())
		return to, nil
	default:
		to := TradeOutcome{}
		return to, errors.New("token number not found")
	}
}

func (p *UniswapV2Pair) PriceImpactToken1BuyToken0(tokenOneBuyAmount *big.Int) (TradeOutcome, *big.Int, *big.Int) {
	to := TradeOutcome{
		AmountIn:            tokenOneBuyAmount,
		AmountInAddr:        p.Token1,
		StartReservesToken0: p.Reserve0,
		StartReservesToken1: p.Reserve1,
	}
	amountInWithFee := new(big.Int).Mul(tokenOneBuyAmount, big.NewInt(997))
	//fmt.Println("amountInWithFee", amountInWithFee.String())
	numerator := new(big.Int).Mul(amountInWithFee, p.Reserve0)
	denominator := new(big.Int).Mul(p.Reserve1, big.NewInt(1000))
	denominator = new(big.Int).Add(denominator, amountInWithFee)
	//fmt.Println("denominator", denominator.String())
	if denominator.Cmp(big.NewInt(0)) == 0 {
		log.Warn().Msg("denominator is 0")
		return to, p.Reserve0, p.Reserve1
	}
	amountOut := new(big.Int).Div(numerator, denominator)
	to.AmountOut = amountOut
	amountInWithFee = new(big.Int).Mul(tokenOneBuyAmount, big.NewInt(3))
	numerator = new(big.Int).Mul(amountInWithFee, p.Reserve0)
	denominator = new(big.Int).Mul(p.Reserve1, big.NewInt(1000))
	denominator = new(big.Int).Add(denominator, amountInWithFee)
	if denominator.Cmp(big.NewInt(0)) == 0 {
		log.Warn().Msg("denominator is 0")
		return to, p.Reserve0, p.Reserve1
	}
	amountOutFee := new(big.Int).Div(numerator, denominator)
	//fmt.Println("amountOut", amountOut.String())
	to.AmountFees = amountOutFee
	p.Reserve1 = new(big.Int).Add(p.Reserve1, tokenOneBuyAmount)
	p.Reserve0 = new(big.Int).Sub(p.Reserve0, amountOut)
	to.EndReservesToken0 = p.Reserve0
	to.EndReservesToken1 = p.Reserve1
	return to, p.Reserve0, p.Reserve1
}

func (p *UniswapV2Pair) PriceImpactToken0BuyToken1(tokenZeroBuyAmount *big.Int) (TradeOutcome, *big.Int, *big.Int) {
	to := TradeOutcome{
		AmountIn:            tokenZeroBuyAmount,
		AmountInAddr:        p.Token0,
		StartReservesToken0: p.Reserve0,
		StartReservesToken1: p.Reserve1,
	}
	amountInWithFee := new(big.Int).Mul(tokenZeroBuyAmount, big.NewInt(997))
	numerator := new(big.Int).Mul(amountInWithFee, p.Reserve1)
	denominator := new(big.Int).Mul(p.Reserve0, big.NewInt(1000))
	denominator = new(big.Int).Add(denominator, amountInWithFee)
	if denominator.Cmp(big.NewInt(0)) == 0 {
		log.Warn().Msg("denominator is 0")
		return to, p.Reserve0, p.Reserve1
	}
	amountOut := new(big.Int).Div(numerator, denominator)
	to.AmountOut = amountOut
	amountInWithFee = new(big.Int).Mul(tokenZeroBuyAmount, big.NewInt(3))
	numerator = new(big.Int).Mul(amountInWithFee, p.Reserve1)
	denominator = new(big.Int).Mul(p.Reserve0, big.NewInt(1000))
	denominator = new(big.Int).Add(denominator, amountInWithFee)
	if denominator.Cmp(big.NewInt(0)) == 0 {
		log.Warn().Msg("denominator is 0")
		return to, p.Reserve0, p.Reserve1
	}
	amountOutFee := new(big.Int).Div(numerator, denominator)
	to.AmountFees = amountOutFee
	p.Reserve0 = new(big.Int).Add(p.Reserve0, tokenZeroBuyAmount)
	p.Reserve1 = new(big.Int).Sub(p.Reserve1, amountOut)
	to.EndReservesToken0 = p.Reserve0
	to.EndReservesToken1 = p.Reserve1
	return to, p.Reserve0, p.Reserve1
}
