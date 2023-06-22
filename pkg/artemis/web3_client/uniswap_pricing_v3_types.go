package web3_client

import (
	"math/big"

	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/entities"
)

type UniswapPoolV3 struct {
	web3_actions.Web3Actions
	*entities.Pool
	PoolAddress          string
	Fee                  constants.FeeAmount
	Slot0                Slot0
	Liquidity            *big.Int
	TokenFeePath         TokenFeePath
	TickListDataProvider *entities.TickListDataProvider
}

type JSONUniswapPoolV3 struct {
	*entities.Pool
	PoolAddress          string
	Fee                  constants.FeeAmount
	Slot0                JSONSlot0
	Liquidity            string
	TokenFeePath         TokenFeePath
	TickListDataProvider *entities.JSONTickListDataProvider
}

func (p *UniswapPoolV3) ConvertToJSONType() JSONUniswapPoolV3 {

	var tickListDataProviderJSON entities.JSONTickListDataProvider
	if p.TickDataProvider != nil {
		tickListDataProviderJSON = p.TickListDataProvider.ConvertToJSONType()
	}
	return JSONUniswapPoolV3{
		PoolAddress:          p.PoolAddress,
		Fee:                  p.Fee,
		Slot0:                p.Slot0.ConvertToJSONType(),
		Liquidity:            p.Liquidity.String(),
		TokenFeePath:         p.TokenFeePath,
		TickListDataProvider: &tickListDataProviderJSON,
	}
}

type Slot0 struct {
	SqrtPriceX96 *big.Int
	Tick         int
	FeeProtocol  int
}

type JSONSlot0 struct {
	SqrtPriceX96 string
	Tick         int
	FeeProtocol  int
}

func (s *JSONSlot0) ConvertToBigIntType() Slot0 {
	price, _ := new(big.Int).SetString(s.SqrtPriceX96, 10)
	return Slot0{
		SqrtPriceX96: price,
		Tick:         s.Tick,
		FeeProtocol:  s.FeeProtocol,
	}
}

func (s *Slot0) ConvertToJSONType() JSONSlot0 {
	return JSONSlot0{
		SqrtPriceX96: s.SqrtPriceX96.String(),
		Tick:         s.Tick,
		FeeProtocol:  s.FeeProtocol,
	}
}
