package artemis_uniswap_pricing

import (
	"errors"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

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
	if p.Token0 == addr {
		return 0
	}
	if p.Token1 == addr {
		return 1
	}
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
