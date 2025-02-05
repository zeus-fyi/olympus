package artemis_uniswap_pricing

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	uniswap_core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/constants"
	entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/utils"
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

func (p *UniswapV3Pair) PriceImpact(ctx context.Context, token *uniswap_core_entities.Token, amountIn *big.Int) (*uniswap_core_entities.CurrencyAmount, *entities.Pool, error) {
	amountInTrade := uniswap_core_entities.FromRawAmount(token, amountIn)
	out, pool, er := p.GetOutputAmount(amountInTrade, nil)
	if er != nil {
		return nil, nil, er
	}
	//adjOut, err := artemis_pricing_utils.ApplyTransferTax(token.Address, out.Quotient())
	//if err != nil {
	//	return nil, nil, err
	//}
	//adjOut = artemis_eth_units.SetSlippage(adjOut)
	//out.Numerator = adjOut
	//out.Denominator = big.NewInt(1)
	return out, pool, nil
}

func (p *UniswapV3Pair) PricingData(ctx context.Context, path artemis_trading_types.TokenFeePath) error {
	// todo, need to handle multi-hops, not sure if this is sufficient for that
	p.Fee = constants.FeeAmount(path.GetFirstFee().Int64())
	wc := p.Web3Actions
	bn, berr := artemis_trading_cache.GetLatestBlockFromCacheOrProvidedSource(context.Background(), wc)
	if berr != nil {
		return berr
	}
	bnst := fmt.Sprintf("%d", bn)
	sessionID := wc.GetSessionLockHeader()
	if wc.GetSessionLockHeader() != "" {
		bnst = fmt.Sprintf("%s-%s", bnst, sessionID)
	}
	hs := crypto.Keccak256Hash([]byte(path.TokenIn.Hex() + bnst + path.GetEndToken().Hex())).String()
	val, ok := Cache.Get(hs)
	if ok && val != nil {
		if assertedVal, tok := val.(*UniswapV3Pair); tok {
			//log.Info().Interface("bn", bn).Interface("pair", assertedVal.PoolAddress).Msg("found v3 pair in cache")
			p.PoolAddress = assertedVal.PoolAddress
			p.Fee = assertedVal.Fee
			p.Slot0 = assertedVal.Slot0
			p.Liquidity = assertedVal.Liquidity
			p.TickListDataProvider = assertedVal.TickListDataProvider
			p.Pool = assertedVal.Pool
			return nil
		} else {
			return fmt.Errorf("value is not of type *entities.Pool")
		}
	}
	pd, er := redisCache.GetPairPricesFromCacheIfExists(ctx, hs)
	if er != nil {
		log.Err(er).Msgf("Error getting v3 pair prices from cache for %s", hs)
		er = nil
	} else {
		assertedVal := pd.V3Pair
		p.PoolAddress = assertedVal.PoolAddress
		p.Fee = assertedVal.Fee
		p.Slot0 = assertedVal.Slot0
		p.Liquidity = assertedVal.Liquidity
		p.TickListDataProvider = assertedVal.TickListDataProvider
		p.Pool = assertedVal.Pool
		return nil
	}

	// todo, store decimals in db
	tokenA := uniswap_core_entities.NewToken(1, accounts.HexToAddress(path.TokenIn.Hex()), 0, "", "")
	tokenB := uniswap_core_entities.NewToken(1, accounts.HexToAddress(path.GetEndToken().Hex()), 0, "", "")
	// todo not sure if this factoryAddress covers all cases
	factoryAddress := accounts.HexToAddress(UniswapV3FactoryAddress)
	pa, err := utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, p.Fee, "")
	if err != nil {
		return err
	}
	p.PoolAddress = pa.String()
	err = p.GetLiquidityAndSlot0FromMulticall3(ctx)
	if err != nil {
		log.Err(err).Msg("PricingData: error getting liquidity and slot0 from multicall3")
		return err
	}
	ts, err := p.GetPopulatedTicksMap()
	if err != nil {
		log.Err(err).Msg("PricingData: error getting populated ticks map")
		return err
	}
	if len(ts) <= 0 {
		return errors.New("PricingData: no populated ticks")
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
	if p != nil && p.Liquidity != nil {
		pricingData := UniswapPricingData{
			V3Pair: *p,
		}
		err = redisCache.AddOrUpdatePairPricesCache(ctx, hs, pricingData, time.Minute*60)
		if err != nil {
			log.Err(err).Msgf("Error adding v3 pair prices to cache for %s", hs)
			err = nil
		} else {
			Cache.Set(hs, p, time.Minute*5)
		}
	}
	return nil
}
