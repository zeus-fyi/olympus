package uniswap_pricing

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/entities"
)

type JSONUniswapV2Pair struct {
	PairContractAddr     string           `json:"pairContractAddr"`
	Price0CumulativeLast *string          `json:"price0CumulativeLast,omitempty"`
	Price1CumulativeLast *string          `json:"price1CumulativeLast,omitempty"`
	KLast                *string          `json:"kLast,omitempty"`
	Token0               accounts.Address `json:"token0"`
	Token1               accounts.Address `json:"token1"`
	Reserve0             string           `json:"reserve0"`
	Reserve1             string           `json:"reserve1"`
	BlockTimestampLast   string           `json:"blockTimestampLast,omitempty"`
}

func (p *JSONUniswapV2Pair) ConvertToBigIntType() *UniswapV2Pair {
	p0 := new(big.Int)
	p1 := new(big.Int)
	k := new(big.Int)
	if p.Price0CumulativeLast != nil {
		p0, _ = new(big.Int).SetString(*p.Price0CumulativeLast, 10)
	}
	if p.Price1CumulativeLast != nil {
		p1, _ = new(big.Int).SetString(*p.Price1CumulativeLast, 10)
	}
	if p.KLast != nil {
		k, _ = new(big.Int).SetString(*p.KLast, 10)
	}

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
	var p0 string
	if p.Price0CumulativeLast != nil {
		p0 = p.Price0CumulativeLast.String()
	}
	var p1 string
	if p.Price1CumulativeLast != nil {
		p1 = p.Price1CumulativeLast.String()
	}
	var k string
	if p.KLast != nil {
		k = p.KLast.String()
	}
	return &JSONUniswapV2Pair{
		PairContractAddr:     p.PairContractAddr,
		Price0CumulativeLast: &p0,
		Price1CumulativeLast: &p1,
		KLast:                &k,
		Token0:               p.Token0,
		Token1:               p.Token1,
		Reserve0:             p.Reserve0.String(),
		Reserve1:             p.Reserve1.String(),
		BlockTimestampLast:   p.BlockTimestampLast.String(),
	}
}

type JSONUniswapPoolV3 struct {
	*entities.Pool       `json:"pool,omitempty"`
	PoolAddress          string                             `json:"poolAddress"`
	Fee                  constants.FeeAmount                `json:"fee"`
	Slot0                JSONSlot0                          `json:"slot0"`
	Liquidity            string                             `json:"liquidity"`
	TokenFeePath         artemis_trading_types.TokenFeePath `json:"tokenFeePath"`
	TickListDataProvider *entities.JSONTickListDataProvider `json:"tickListDataProvider,omitempty"`
}

func (p *JSONUniswapPoolV3) ConvertToBigIntType() *UniswapV3Pair {
	lq, _ := new(big.Int).SetString(p.Liquidity, 10)
	var tl *entities.JSONTickListDataProvider
	var tlBigInt *entities.TickListDataProvider
	if p.TickListDataProvider != nil {
		tl = p.TickListDataProvider
		tlBigInt = tl.ConvertToBigIntType()
	}
	return &UniswapV3Pair{
		Pool:                 p.Pool,
		PoolAddress:          p.PoolAddress,
		Fee:                  p.Fee,
		Slot0:                p.Slot0.ConvertToBigIntType(),
		Liquidity:            lq,
		TokenFeePath:         p.TokenFeePath,
		TickListDataProvider: tlBigInt,
	}
}

func (p *UniswapV3Pair) ConvertToJSONType() *JSONUniswapPoolV3 {
	var tickListDataProviderJSON entities.JSONTickListDataProvider
	if p.TickDataProvider != nil {
		tickListDataProviderJSON = p.TickListDataProvider.ConvertToJSONType()
	}
	return &JSONUniswapPoolV3{
		PoolAddress:          p.PoolAddress,
		Fee:                  p.Fee,
		Slot0:                p.Slot0.ConvertToJSONType(),
		Liquidity:            p.Liquidity.String(),
		TokenFeePath:         p.TokenFeePath,
		TickListDataProvider: &tickListDataProviderJSON,
	}
}
