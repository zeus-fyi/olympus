package web3_client

import (
	"context"
	"errors"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
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
	adjOut := ApplyTransferTax(token.Address, out.Quotient())
	out.Numerator = adjOut
	out.Denominator = big.NewInt(1)
	return out, pool, nil
}

/*

type Erc20TokenInfo struct {
	Address                string  `db:"address" json:"address"`
	ProtocolNetworkID      int     `db:"protocol_network_id" json:"protocolNetworkID"`
	BalanceOfSlotNum       int     `db:"balance_of_slot_num" json:"balanceOfSlotNum"`
	Name                   *string `db:"name" json:"name"`
	Symbol                 *string `db:"symbol" json:"symbol"`
	Decimals               *int    `db:"decimals" json:"decimals"`
	TransferTaxNumerator   *int    `db:"transfer_tax_numerator" json:"transferTaxNumerator"`
	TransferTaxDenominator *int    `db:"transfer_tax_denominator" json:"transferTaxDenominator"`
	TradingEnabled         *bool   `db:"trading_enabled" json:"tradingEnabled"`
}
*/

func (p *UniswapPoolV3) PricingData(ctx context.Context, path TokenFeePath, simMode bool) error {
	// todo, need to handle multi-hops, not sure if this is sufficient for that
	p.Fee = constants.FeeAmount(path.GetFirstFee().Int64())
	p.SimMode = simMode
	wc := p.Web3Actions
	if artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeHistoricalData.NodeURL != "" && !simMode {
		wc = web3_actions.NewWeb3ActionsClientWithAccount(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeHistoricalData.NodeURL, artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeHistoricalData.Account)
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
		tokenCurrencyA := core_entities.NewToken(1, accounts.HexToAddress(tokenA.Address), uint(decimals), sym, name)
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
		tokenCurrencyB := core_entities.NewToken(1, accounts.HexToAddress(path.GetEndToken().Hex()), uint(decimals), sym, name)
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
		tokenA := core_entities.NewToken(1, accounts.HexToAddress(path.TokenIn.Hex()), uint(decimals), "", "")
		decimals, err = wc.GetContractDecimals(ctx, path.GetEndToken().Hex())
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
	return nil
}
