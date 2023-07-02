package uniswap_pricing

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_pricing_utils "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/utils"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
)

const (
	getReserves          = "getReserves"
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
	Price0CumulativeLast *big.Int         `json:"price0CumulativeLast"`
	Price1CumulativeLast *big.Int         `json:"price1CumulativeLast"`
	KLast                *big.Int         `json:"kLast"`
	Token0               accounts.Address `json:"token0"`
	Token1               accounts.Address `json:"token1"`
	Reserve0             *big.Int         `json:"reserve0"`
	Reserve1             *big.Int         `json:"reserve1"`
	BlockTimestampLast   *big.Int         `json:"blockTimestampLast"`
}

func (p *UniswapV2Pair) PairForV2FromAddresses(tokenA, tokenB accounts.Address) error {
	return p.PairForV2(tokenA.String(), tokenB.String())
}

func (p *UniswapV2Pair) PairForV2(tokenA, tokenB string) error {
	if tokenA == ZeroEthAddress {
		tokenA = WETH.String()
	}
	if tokenB == ZeroEthAddress {
		tokenB = WETH.String()
	}
	if tokenA == tokenB {
		return errors.New("identical addresses")
	}
	p.sortTokens(accounts.HexToAddress(tokenA), accounts.HexToAddress(tokenB))
	message := []byte{255}
	message = append(message, common.HexToAddress(UniswapV2FactoryAddress.String()).Bytes()...)
	addrSum := p.Token0.Bytes()
	addrSum = append(addrSum, p.Token1.Bytes()...)
	message = append(message, crypto.Keccak256(addrSum)...)
	b, err := hex.DecodeString(pairAddressSuffix)
	if err != nil {
		return err
	}
	message = append(message, b...)
	hashed := crypto.Keccak256(message)
	addressBytes := big.NewInt(0).SetBytes(hashed)
	addressBytes = addressBytes.Abs(addressBytes)
	p.PairContractAddr = common.BytesToAddress(addressBytes.Bytes()).String()
	return nil
}

func (p *UniswapV2Pair) GetQuoteUsingTokenAddr(addr string, amount *big.Int) (*big.Int, error) {
	if addr == "0x0000000000000000000000000000000000000000" {
		addr = WETH.String()
	}
	if p.Token0 == accounts.HexToAddress(addr) {
		return p.GetQuoteToken0BuyToken1(amount)
	}
	if p.Token1 == accounts.HexToAddress(addr) {
		return p.GetQuoteToken1BuyToken0(amount)
	}
	return nil, errors.New("GetQuoteUsingTokenAddr: token not found")
}

func (p *UniswapV2Pair) GetOppositeToken(addr string) accounts.Address {
	if addr == "0x0000000000000000000000000000000000000000" {
		addr = WETH.String()
	}
	if p.Token0 == accounts.HexToAddress(addr) {
		return p.Token1
	}
	if p.Token1 == accounts.HexToAddress(addr) {
		return p.Token0
	}
	log.Warn().Msgf("GetOppositeToken: token not found: %s", addr)
	return accounts.Address{}
}

func (p *UniswapV2Pair) GetTokenNumber(addr accounts.Address) int {
	if addr.String() == "0x0000000000000000000000000000000000000000" {
		if p.Token0.String() == WETH.String() {
			return 0
		}
		if p.Token1.String() == WETH.String() {
			return 1
		}
	}
	if p.Token0 == addr {
		return 0
	}
	if p.Token1 == addr {
		return 1
	}
	log.Warn().Msgf("GetTokenNumber: token not found: %s", addr)
	return -1
}

func (p *UniswapV2Pair) sortTokens(tkn0, tkn1 accounts.Address) {
	token0Rep := big.NewInt(0).SetBytes(tkn0.Bytes())
	token1Rep := big.NewInt(0).SetBytes(tkn1.Bytes())

	if token0Rep.Cmp(token1Rep) > 0 {
		tkn0, tkn1 = tkn1, tkn0
	}
	p.Token0 = tkn0
	p.Token1 = tkn1
}

/*
	price0, err := u.SingleReadMethodBigInt(ctx, "price0CumulativeLast", scInfo)
	if err != nil {
		return uniswap_pricing.UniswapV2Pair{}, err
	}
	pairInfo.Price0CumulativeLast = price0
	price1, err := u.SingleReadMethodBigInt(ctx, "price1CumulativeLast", scInfo)
	if err != nil {
		return uniswap_pricing.UniswapV2Pair{}, err
	}
	pairInfo.Price1CumulativeLast = price1
	kLast, err := u.SingleReadMethodBigInt(ctx, "kLast", scInfo)
	if err != nil {
		return uniswap_pricing.UniswapV2Pair{}, err
	}
	pairInfo.KLast = kLast
	token0, err := u.SingleReadMethodAddr(ctx, "token0", scInfo)
	if err != nil {
		return uniswap_pricing.UniswapV2Pair{}, err
	}
	pairInfo.Token0 = token0
	token1, err := u.SingleReadMethodAddr(ctx, "token1", scInfo)
	if err != nil {
		return uniswap_pricing.UniswapV2Pair{}, err
	}
*/

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

func (p *UniswapV2Pair) PriceImpact(tokenAddrPath accounts.Address, tokenBuyAmount *big.Int) (artemis_trading_types.TradeOutcome, error) {
	tokenNumber := p.GetTokenNumber(tokenAddrPath)
	switch tokenNumber {
	case 1:
		to, _, _ := p.PriceImpactToken1BuyToken0(tokenBuyAmount)
		to.AmountInAddr = tokenAddrPath
		to.AmountOutAddr = p.GetOppositeToken(tokenAddrPath.String())
		to.AmountOut = artemis_pricing_utils.ApplyTransferTax(to.AmountOutAddr, to.AmountOut)
		return to, nil
	case 0:
		to, _, _ := p.PriceImpactToken0BuyToken1(tokenBuyAmount)
		to.AmountInAddr = tokenAddrPath
		to.AmountOutAddr = p.GetOppositeToken(tokenAddrPath.String())
		to.AmountOut = artemis_pricing_utils.ApplyTransferTax(to.AmountOutAddr, to.AmountOut)
		return to, nil
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
