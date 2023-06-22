package web3_client

import (
	"context"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
	uniswap_core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/utils"
)

const (
	ticks                   = "ticks"
	slot0                   = "slot0"
	liquidity               = "liquidity"
	tickBitmap              = "tickBitmap"
	getPopulatedTicksInWord = "getPopulatedTicksInWord"

	TickLensAddress         = "0xbfd8137f7d1516D3ea5cA83523914859ec47F573"
	UniswapV3FactoryAddress = "0x1F98431c8aD98523631AE4a59f267346ea31F984"
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
	// todo, need to handle multi-hops, not sure if this is sufficient for that
	p.Fee = constants.FeeAmount(path.GetFirstFee().Int64())
	decimals, err := p.GetContractDecimals(ctx, path.TokenIn.Hex())
	if err != nil {
		return err
	}
	// todo, store decimals in db
	tokenA := core_entities.NewToken(1, accounts.HexToAddress(path.TokenIn.Hex()), uint(decimals), "", "")
	decimals, err = p.GetContractDecimals(ctx, path.GetEndToken().Hex())
	if err != nil {
		return err
	}
	tokenB := core_entities.NewToken(1, accounts.HexToAddress(path.GetEndToken().Hex()), uint(decimals), "", "")
	// todo not sure if this factoryAddress covers all cases
	factoryAddress := accounts.HexToAddress(UniswapV3FactoryAddress)
	pa, err := utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, p.Fee, "")
	if err != nil {
		return err
	}
	p.PoolAddress = pa.String()
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
