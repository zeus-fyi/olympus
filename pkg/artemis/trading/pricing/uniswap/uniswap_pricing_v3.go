package artemis_uniswap_pricing

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
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
	out, pool, err := p.GetOutputAmount(amountInTrade, nil)
	if err != nil {
		return nil, nil, err
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
	bn, berr := artemis_trading_cache.GetLatestBlockFromCacheOrProvidedSource(ctx, wc)
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
	tm := artemis_trading_cache.TokenMap
	if tm != nil && tm[path.TokenIn.String()].TransferTaxDenominator != nil && tm[path.GetEndToken().String()].TransferTaxDenominator != nil {
		decimals := 0
		tokenA := tm[path.TokenIn.String()]
		if tokenA.Decimals != nil {
			decimals = *tokenA.Decimals
		}
		sym := ""
		if tokenA.Symbol != nil {
			sym = *tokenA.Symbol
		}
		name := ""
		if tokenA.Name != nil {
			name = *tokenA.Name
		}
		tokenCurrencyA := uniswap_core_entities.NewToken(1, accounts.HexToAddress(tokenA.Address), uint(decimals), sym, name)
		tokenB := tm[path.GetEndToken().String()]
		if tokenB.Decimals != nil {
			decimals = *tokenB.Decimals
		}
		sym = ""
		if tokenB.Symbol != nil {
			sym = *tokenB.Symbol
		}
		name = ""
		if tokenB.Name != nil {
			name = *tokenB.Name
		}
		tokenCurrencyB := uniswap_core_entities.NewToken(1, accounts.HexToAddress(path.GetEndToken().Hex()), uint(decimals), sym, name)
		factoryAddress := accounts.HexToAddress(UniswapV3FactoryAddress)
		pa, err := utils.ComputePoolAddress(factoryAddress, tokenCurrencyA, tokenCurrencyB, p.Fee, "")
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
		if len(ts) <= 0 {
			return errors.New("no populated ticks")
		}
		tdp, err := entities.NewTickListDataProvider(ts, constants.TickSpacings[p.Fee])
		if err != nil {
			return err
		}
		p.TickListDataProvider = tdp
		v3Pool, err := entities.NewPool(tokenCurrencyA, tokenCurrencyB, p.Fee, p.Slot0.SqrtPriceX96, p.Liquidity, p.Slot0.Tick, p.TickListDataProvider)
		if err != nil {
			return err
		}
		p.Pool = v3Pool
	} else {
		decimals, err := wc.GetContractDecimals(ctx, path.TokenIn.Hex())
		if err != nil {
			return err
		}
		// todo, store decimals in db
		tokenA := uniswap_core_entities.NewToken(1, accounts.HexToAddress(path.TokenIn.Hex()), uint(decimals), "", "")
		decimals, err = wc.GetContractDecimals(ctx, path.GetEndToken().Hex())
		if err != nil {
			return err
		}
		tokenB := uniswap_core_entities.NewToken(1, accounts.HexToAddress(path.GetEndToken().Hex()), uint(decimals), "", "")
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
		if len(ts) <= 0 {
			return errors.New("no populated ticks")
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
	}
	Cache.Set(hs, p, time.Minute*5)
	return nil
}
