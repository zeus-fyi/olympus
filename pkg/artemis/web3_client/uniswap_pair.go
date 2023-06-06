package web3_client

import (
	"context"
	"errors"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
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

func (p *UniswapV2Pair) GetQuoteUsingTokenAddr(addr string, amount *big.Int) (*big.Int, error) {
	if addr == "0x0000000000000000000000000000000000000000" {
		addr = WETH9ContractAddress
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
		addr = WETH9ContractAddress
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
		if p.Token0.String() == WETH9ContractAddress {
			return 0
		}
		if p.Token1.String() == WETH9ContractAddress {
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

func (u *UniswapClient) GetPairContractPrices(ctx context.Context, pairContractAddr string) (UniswapV2Pair, error) {
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
		Token0:               accounts.Address{},
		Token1:               accounts.Address{},
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

func (u *UniswapClient) SingleReadMethodBigInt(ctx context.Context, methodName string, scInfo *web3_actions.SendContractTxPayload) (*big.Int, error) {
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

func (u *UniswapClient) SingleReadMethodAddr(ctx context.Context, methodName string, scInfo *web3_actions.SendContractTxPayload) (accounts.Address, error) {
	scInfo.MethodName = methodName
	resp, err := u.Web3Client.CallConstantFunction(ctx, scInfo)
	if err != nil {
		return accounts.Address{}, err
	}
	if len(resp) == 0 {
		return accounts.Address{}, errors.New("empty response")
	}
	addr, err := ConvertToAddress(resp[0])
	if err != nil {
		return accounts.Address{}, err
	}
	return addr, nil
}
