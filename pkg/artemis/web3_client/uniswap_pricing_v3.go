package web3_client

import (
	"context"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
	uniswap_core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/entities"
)

const (
	ticks                   = "ticks"
	slot0                   = "slot0"
	liquidity               = "liquidity"
	tickBitmap              = "tickBitmap"
	getPopulatedTicksInWord = "getPopulatedTicksInWord"

	TickLensAddress = "0xbfd8137f7d1516D3ea5cA83523914859ec47F573"
)

func (p *UniswapPoolV3) PriceImpact(ctx context.Context, token *core_entities.Token, amountIn *big.Int) (*uniswap_core_entities.CurrencyAmount, *entities.Pool, error) {
	amountInTrade := uniswap_core_entities.FromRawAmount(token, amountIn)
	out, pool, err := p.GetOutputAmount(amountInTrade, nil)
	if err != nil {
		return nil, nil, err
	}
	return out, pool, nil
}

func (p *UniswapPoolV3) PricingData(ctx context.Context, path TokenFeePath) error {
	decimals, err := p.GetContractDecimals(ctx, path.TokenIn.Hex())
	if err != nil {
		return err
	}
	tokenA := core_entities.NewToken(1, accounts.HexToAddress(path.TokenIn.Hex()), uint(decimals), "", "")
	decimals, err = p.GetContractDecimals(ctx, path.GetEndToken().Hex())
	if err != nil {
		return err
	}
	tokenB := core_entities.NewToken(1, accounts.HexToAddress(path.GetEndToken().Hex()), uint(decimals), "", "")
	err = p.GetSlot0(ctx)
	if err != nil {
		return err
	}
	err = p.GetLiquidity(ctx)
	if err != nil {
		return err
	}
	ts, err := p.GetPopulatedTicksMap()
	if err != nil {
		return err
	}
	// todo get fee from pool vs hardcode
	if p.Fee == 0 {
		p.Fee = constants.FeeMedium
	}
	tdp, err := entities.NewTickListDataProvider(ts, constants.TickSpacings[p.Fee])
	if err != nil {
		return err
	}
	p.TickListDataProvider = tdp
	v3Pool, err := entities.NewPool(tokenA, tokenB, p.Fee, p.Slot0.SqrtPriceX96, p.Liquidity, p.Slot0.Tick, p.TickListDataProvider)
	if err != nil {
		return err
	}
	p.Pool = v3Pool
	return nil
}
