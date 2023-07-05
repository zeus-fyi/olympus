package web3_client

import (
	"fmt"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client/uniswap_libs/uniswap_v3/utils"
)

// example v3 pool: 0x4b5Ab61593A2401B1075b90c04cBCDD3F87CE011

func (s *Web3ClientTestSuite) TestUniswapV3DataFetcherV2() {
	factoryAddress := accounts.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984")
	tokenA := core_entities.NewToken(1, accounts.HexToAddress(WETH9ContractAddress), 18, "WETH", "Wrapped Ether")
	tokenB := core_entities.NewToken(1, accounts.HexToAddress(UsdCoinAddr), 6, "USDC", "USD Coin")
	result, err := utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeMedium, "")
	s.Require().NoError(err)

	p := uniswap_pricing.UniswapV3Pair{
		PoolAddress: result.String(),
		Web3Actions: s.MainnetWeb3User.Web3Actions,
		Fee:         constants.FeeMedium,
	}
	tfp := artemis_trading_types.TokenFeePath{
		TokenIn: tokenA.Address,
		Path: []artemis_trading_types.TokenFee{
			{
				Token: tokenB.Address,
				Fee:   new(big.Int).SetInt64(int64(constants.FeeMedium)),
			},
		},
	}

	err = p.PricingData(ctx, tfp, false)
	s.Require().NoError(err)
	fmt.Println(p.PoolAddress)
	output, _, err := p.PriceImpact(ctx, tokenA, Ether)
	s.Require().NoError(err)
	s.Require().NotNil(output)
	usdAmountSim := new(big.Int).Div(output.Numerator, new(big.Int).SetInt64(1000000))
	fmt.Println("usdAmountSim", usdAmountSim.String())
}

func (s *Web3ClientTestSuite) TestUniswapV3DataFetcher() {
	factoryAddress := accounts.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984")
	tokenA := core_entities.NewToken(1, accounts.HexToAddress(WETH9ContractAddress), 18, "WETH", "Wrapped Ether")
	tokenB := core_entities.NewToken(1, accounts.HexToAddress(UsdCoinAddr), 6, "USDC", "USD Coin")
	result, err := utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeMedium, "")
	s.Require().NoError(err)
	p := uniswap_pricing.UniswapV3Pair{
		PoolAddress: result.String(),
		Web3Actions: s.MainnetWeb3User.Web3Actions,
		Fee:         constants.FeeMedium,
	}
	fmt.Println(p.PoolAddress)
	err = p.GetSlot0(ctx)
	s.Require().NoError(err)

	err = p.GetLiquidity(ctx)
	s.Require().NoError(err)
	ts, err := p.GetPopulatedTicksMap()
	s.Require().NoError(err)
	s.Require().NotEmpty(ts)

	tdp, err := entities.NewTickListDataProvider(ts, constants.TickSpacings[constants.FeeMedium])
	s.Require().NoError(err)

	v3Pool, err := entities.NewPool(tokenA, tokenB, constants.FeeMedium, p.Slot0.SqrtPriceX96, p.Liquidity, p.Slot0.Tick, tdp)
	s.Require().NoError(err)
	s.Require().NotNil(v3Pool)

	inputAmount := core_entities.FromRawAmount(tokenA, Ether)

	output, pool, err := v3Pool.GetOutputAmount(inputAmount, nil)
	s.Require().NoError(err)
	s.Require().NotNil(output)
	s.Require().NotNil(pool)

	usdAmountSim := new(big.Int).Div(output.Numerator, new(big.Int).SetInt64(1000000))
	fmt.Println("usdAmountSim", usdAmountSim.String())

	uni := InitUniswapClient(ctx, s.MainnetWeb3User)
	tfp := artemis_trading_types.TokenFeePath{
		TokenIn: tokenA.Address,
		Path: []artemis_trading_types.TokenFee{{
			Token: tokenB.Address,
			Fee:   new(big.Int).SetInt64(int64(v3Pool.Fee)),
		}},
	}
	qp := QuoteExactInputSingleParams{
		TokenIn:           tfp.TokenIn,
		TokenOut:          tfp.GetEndToken(),
		Fee:               new(big.Int).SetInt64(int64(v3Pool.Fee)),
		AmountIn:          Ether,
		SqrtPriceLimitX96: big.NewInt(0),
	}
	resp, err := uni.GetPoolV3ExactInputSingleQuoteFromQuoterV2(ctx, qp)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	usdAmountActual := new(big.Int).Div(resp.AmountOut, new(big.Int).SetInt64(1000000))
	fmt.Println("usdcAmount", usdAmountActual.String())

	s.Equal(usdAmountActual.String(), usdAmountSim.String())
}

func (s *Web3ClientTestSuite) TestUniswapV3() {
	factoryAddress := accounts.HexToAddress("0x1111111111111111111111111111111111111111")
	tokenA := core_entities.NewToken(1, accounts.HexToAddress(UsdCoinAddr), 18, "USDC", "USD Coin")
	tokenB := core_entities.NewToken(1, accounts.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"), 18, "DAI", "Dai Stablecoin")

	result, err := utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeLow, "")
	if err != nil {
		panic(err)
	}
	s.Equal(result, accounts.HexToAddress("0x90B1b09A9715CaDbFD9331b3A7652B24BfBEfD32"))

	v3Pool, err := entities.NewPool(tokenA, tokenB, constants.FeeLow, nil, nil, 0, nil)
	s.Require().NoError(err)
	s.Require().NotNil(v3Pool)

	inputAmount := core_entities.FromRawAmount(tokenA, big.NewInt(100))
	output, pool, err := v3Pool.GetOutputAmount(inputAmount, nil)
	s.Require().NoError(err)
	s.Require().NotNil(output)
	s.Require().NotNil(pool)
}
