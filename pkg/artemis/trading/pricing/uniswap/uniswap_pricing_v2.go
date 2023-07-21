package artemis_uniswap_pricing

import (
	"errors"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/constants"
)

const (
	pairAddressSuffix    = "96e8ac4277198ff8b6f785478aa9a39f403cb768dd02cbee326c3e7da348845f"
	ZeroEthAddress       = "0x0000000000000000000000000000000000000000"
	WETH9ContractAddress = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
)

var (
	WETH                    = accounts.HexToAddress(WETH9ContractAddress)
	UniswapV2FactoryAddress = accounts.HexToAddress("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f")
)

type UniswapV2Pair struct {
	PairContractAddr     string           `json:"pairContractAddr"`
	Price0CumulativeLast *big.Int         `json:"price0CumulativeLast,omitempty"`
	Price1CumulativeLast *big.Int         `json:"price1CumulativeLast,omitempty"`
	KLast                *big.Int         `json:"kLast,omitempty"`
	Token0               accounts.Address `json:"token0"`
	Token1               accounts.Address `json:"token1"`
	Reserve0             *big.Int         `json:"reserve0,"`
	Reserve1             *big.Int         `json:"reserve1"`
	BlockTimestampLast   *big.Int         `json:"blockTimestampLast,omitempty"`
}

func (p *UniswapV2Pair) GetBaseFee() constants.FeeAmount {
	return constants.FeeMedium
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

func (p *UniswapV2Pair) PriceImpact(tokenAddrPath accounts.Address, tokenBuyAmount *big.Int) (artemis_trading_types.TradeOutcome, error) {
	tokenNumber := p.GetTokenNumber(tokenAddrPath)
	var err error

	switch tokenNumber {
	case 1:
		to, _, _ := p.PriceImpactToken1BuyToken0(tokenBuyAmount)
		to.AmountInAddr = tokenAddrPath
		to.AmountOutAddr = p.GetOppositeToken(tokenAddrPath.String())
		//to.AmountOut, err = artemis_pricing_utils.ApplyTransferTax(to.AmountOutAddr, to.AmountOut)
		//if err != nil {
		//	return to, err
		//}
		//to.AmountOut, err = artemis_pricing_utils.ApplyTransferTax(to.AmountInAddr, to.AmountOut)
		//if err != nil {
		//	return to, err
		//}
		return to, err
	case 0:
		to, _, _ := p.PriceImpactToken0BuyToken1(tokenBuyAmount)
		to.AmountInAddr = tokenAddrPath
		to.AmountOutAddr = p.GetOppositeToken(tokenAddrPath.String())
		//to.AmountOut, err = artemis_pricing_utils.ApplyTransferTax(to.AmountOutAddr, to.AmountOut)
		//if err != nil {
		//	return to, err
		//}
		//to.AmountOut, err = artemis_pricing_utils.ApplyTransferTax(to.AmountInAddr, to.AmountOut)
		//if err != nil {
		//	return to, err
		//}
		return to, err
	default:
		to := artemis_trading_types.TradeOutcome{}
		return to, errors.New("token number not found")
	}
}

func (p *UniswapV2Pair) PriceImpactToken1BuyToken0(tokenOneBuyAmount *big.Int) (artemis_trading_types.TradeOutcome, *big.Int, *big.Int) {
	to := artemis_trading_types.TradeOutcome{
		AmountIn:            tokenOneBuyAmount,
		AmountInAddr:        p.Token1,
		StartReservesToken0: p.Reserve0,
		StartReservesToken1: p.Reserve1,
	}
	if tokenOneBuyAmount == nil {
		tokenOneBuyAmount = big.NewInt(0)
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

func (p *UniswapV2Pair) PriceImpactToken0BuyToken1(tokenZeroBuyAmount *big.Int) (artemis_trading_types.TradeOutcome, *big.Int, *big.Int) {
	to := artemis_trading_types.TradeOutcome{
		AmountIn:            tokenZeroBuyAmount,
		AmountInAddr:        p.Token0,
		StartReservesToken0: p.Reserve0,
		StartReservesToken1: p.Reserve1,
	}
	if tokenZeroBuyAmount == nil {
		tokenZeroBuyAmount = big.NewInt(0)
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
