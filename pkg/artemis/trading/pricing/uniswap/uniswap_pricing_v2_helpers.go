package uniswap_pricing

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

type JSONUniswapV2Pair struct {
	PairContractAddr     string           `json:"pairContractAddr"`
	Price0CumulativeLast string           `json:"price0CumulativeLast,omitempty"`
	Price1CumulativeLast string           `json:"price1CumulativeLast,omitempty"`
	KLast                string           `json:"kLast,omitempty"`
	Token0               accounts.Address `json:"token0"`
	Token1               accounts.Address `json:"token1"`
	Reserve0             string           `json:"reserve0"`
	Reserve1             string           `json:"reserve1"`
	BlockTimestampLast   string           `json:"blockTimestampLast,omitempty"`
}

func (p *JSONUniswapV2Pair) ConvertToBigIntType() *UniswapV2Pair {
	p0, _ := new(big.Int).SetString(p.Price0CumulativeLast, 10)
	p1, _ := new(big.Int).SetString(p.Price1CumulativeLast, 10)
	k, _ := new(big.Int).SetString(p.KLast, 10)
	r0, _ := new(big.Int).SetString(p.Reserve0, 10)
	r1, _ := new(big.Int).SetString(p.Reserve1, 10)
	bt, _ := new(big.Int).SetString(p.BlockTimestampLast, 10)
	return &UniswapV2Pair{
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
func (p *UniswapV2Pair) ConvertToJSONType() *JSONUniswapV2Pair {
	if p.Price0CumulativeLast == nil {
		p.Price0CumulativeLast = big.NewInt(0)
	}
	if p.Price1CumulativeLast == nil {
		p.Price1CumulativeLast = big.NewInt(0)
	}
	if p.KLast == nil {
		p.KLast = big.NewInt(0)

	}
	return &JSONUniswapV2Pair{
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
