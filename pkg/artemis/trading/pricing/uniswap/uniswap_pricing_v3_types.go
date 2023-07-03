package uniswap_pricing

import (
	"math/big"

	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/entities"
)

type UniswapPoolV3 struct {
	web3_actions.Web3Actions `json:"-,omitempty"`
	*entities.Pool           `json:"pool,omitempty"`
	PoolAddress              string                             `json:"poolAddress"`
	Fee                      constants.FeeAmount                `json:"fee"`
	Slot0                    Slot0                              `json:"slot0"`
	Liquidity                *big.Int                           `json:"liquidity"`
	TokenFeePath             artemis_trading_types.TokenFeePath `json:"tokenFeePath"`
	TickListDataProvider     *entities.TickListDataProvider     `json:"tickListDataProvider,omitempty"`
	SimMode                  bool                               `json:"simMode,omitempty"`
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
	SqrtPriceX96 *big.Int `json:"sqrtPriceX96"`
	Tick         int      `json:"tick"`
	FeeProtocol  int      `json:"feeProtocol"`
}

type JSONSlot0 struct {
	SqrtPriceX96 string `json:"sqrtPriceX96"`
	Tick         int    `json:"tick"`
	FeeProtocol  int    `json:"feeProtocol"`
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
