package web3_client

import (
	"errors"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

type JSONUniswapV2Pair struct {
	PairContractAddr     string           `json:"pairContractAddr"`
	Price0CumulativeLast string           `json:"price0CumulativeLast"`
	Price1CumulativeLast string           `json:"price1CumulativeLast"`
	KLast                string           `json:"kLast"`
	Token0               accounts.Address `json:"token0"`
	Token1               accounts.Address `json:"token1"`
	Reserve0             string           `json:"reserve0"`
	Reserve1             string           `json:"reserve1"`
	BlockTimestampLast   string           `json:"blockTimestampLast"`
}

func (p *JSONUniswapV2Pair) ConvertToBigIntType() UniswapV2Pair {
	p0, _ := new(big.Int).SetString(p.Price0CumulativeLast, 10)
	p1, _ := new(big.Int).SetString(p.Price1CumulativeLast, 10)
	k, _ := new(big.Int).SetString(p.KLast, 10)
	r0, _ := new(big.Int).SetString(p.Reserve0, 10)
	r1, _ := new(big.Int).SetString(p.Reserve1, 10)
	bt, _ := new(big.Int).SetString(p.BlockTimestampLast, 10)
	return UniswapV2Pair{
		PairContractAddr:     p.PairContractAddr,
		Price0CumulativeLast: p0,
		Price1CumulativeLast: p1,
		KLast:                k,
		Token0:               p.Token0,
		Token1:               p.Token1,
		Reserve0:             r0,
		Reserve1:             r1,
		BlockTimestampLast:   bt,
	}
}
func (p *UniswapV2Pair) ConvertToJSONType() JSONUniswapV2Pair {
	return JSONUniswapV2Pair{
		PairContractAddr:     p.PairContractAddr,
		Price0CumulativeLast: p.Price0CumulativeLast.String(),
		Price1CumulativeLast: p.Price1CumulativeLast.String(),
		KLast:                p.KLast.String(),
		Token0:               p.Token0,
		Token1:               p.Token1,
		Reserve0:             p.Reserve0.String(),
		Reserve1:             p.Reserve1.String(),
		BlockTimestampLast:   p.BlockTimestampLast.String(),
	}
}

func (p *UniswapV2Pair) GetQuoteToken0BuyToken1(token0 *big.Int) (*big.Int, error) {
	if p.Reserve0 == nil || p.Reserve1 == nil || p.Reserve0.Cmp(big.NewInt(0)) == 0 || p.Reserve1.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("reserves are not initialized or are zero")
	}
	amountInWithFee := new(big.Int).Mul(token0, big.NewInt(997))
	numerator := new(big.Int).Mul(amountInWithFee, p.Reserve1)
	denominator := new(big.Int).Mul(p.Reserve0, big.NewInt(1000))
	denominator = new(big.Int).Add(denominator, amountInWithFee)
	if denominator.Cmp(big.NewInt(0)) == 0 {
		log.Warn().Msg("denominator is 0")
		return nil, errors.New("denominator is 0")
	}
	amountOut := new(big.Int).Div(numerator, denominator)
	return amountOut, nil
}

func (p *UniswapV2Pair) GetQuoteToken1BuyToken0(token1 *big.Int) (*big.Int, error) {
	if p.Reserve0 == nil || p.Reserve1 == nil || p.Reserve0.Cmp(big.NewInt(0)) == 0 || p.Reserve1.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("reserves are not initialized or are zero")
	}
	amountInWithFee := new(big.Int).Mul(token1, big.NewInt(997))
	numerator := new(big.Int).Mul(amountInWithFee, p.Reserve0)
	denominator := new(big.Int).Mul(p.Reserve1, big.NewInt(1000))
	denominator = new(big.Int).Add(denominator, amountInWithFee)
	if denominator.Cmp(big.NewInt(0)) == 0 {
		log.Warn().Msg("denominator is 0")
		return nil, errors.New("denominator is 0")
	}
	amountOut := new(big.Int).Div(numerator, denominator)
	return amountOut, nil
}
