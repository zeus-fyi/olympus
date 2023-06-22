package periphery

import (
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	uniswap_core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/utils"
)

var (
	ether  = uniswap_core_entities.EtherOnChain(1)
	token0 = uniswap_core_entities.NewToken(1, accounts.HexToAddress("0x0000000000000000000000000000000000000001"), 18, "t0", "token0")
	token1 = uniswap_core_entities.NewToken(1, accounts.HexToAddress("0x0000000000000000000000000000000000000002"), 18, "t1", "token1")
	token2 = uniswap_core_entities.NewToken(1, accounts.HexToAddress("0x0000000000000000000000000000000000000003"), 18, "t2", "token2")
	token3 = uniswap_core_entities.NewToken(1, accounts.HexToAddress("0x0000000000000000000000000000000000000004"), 18, "t2", "token3")

	weth = ether.Wrapped()

	pool_0_1_medium, _ = entities.NewPool(token0, token1, constants.FeeMedium, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	pool_1_2_low, _    = entities.NewPool(token1, token2, constants.FeeLow, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	pool_0_weth, _     = entities.NewPool(token0, weth, constants.FeeMedium, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)
	pool_1_weth, _     = entities.NewPool(token1, weth, constants.FeeMedium, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)

	route_0_1, _   = entities.NewRoute([]*entities.Pool{pool_0_1_medium}, token0, token1)
	route_0_1_2, _ = entities.NewRoute([]*entities.Pool{pool_0_1_medium, pool_1_2_low}, token0, token2)

	route_0_weth, _   = entities.NewRoute([]*entities.Pool{pool_0_weth}, token0, weth)
	route_0_1_weth, _ = entities.NewRoute([]*entities.Pool{pool_0_1_medium, pool_1_weth}, token0, weth)
	route_weth_0, _   = entities.NewRoute([]*entities.Pool{pool_0_weth}, weth, token0)
	route_weth_0_1, _ = entities.NewRoute([]*entities.Pool{pool_0_weth, pool_0_1_medium}, weth, token1)

	feeAmount    = constants.FeeMedium
	sqrtRatioX96 = utils.EncodeSqrtRatioX96(big.NewInt(1), big.NewInt(1))
	liquidity    = big.NewInt(1_000_000)
	tick, _      = utils.GetTickAtSqrtRatio(sqrtRatioX96)
	ticks        = []entities.Tick{
		{
			Index:          entities.NearestUsableTick(utils.MinTick, constants.TickSpacings[feeAmount]),
			LiquidityNet:   liquidity,
			LiquidityGross: liquidity,
		},
		{
			Index:          entities.NearestUsableTick(utils.MaxTick, constants.TickSpacings[feeAmount]),
			LiquidityNet:   new(big.Int).Mul(liquidity, constants.NegativeOne),
			LiquidityGross: liquidity,
		},
	}

	p, _     = entities.NewTickListDataProvider(ticks, constants.TickSpacings[feeAmount])
	makePool = func(token0, token1 *uniswap_core_entities.Token) *entities.Pool {
		pool, _ := entities.NewPool(token0, token1, feeAmount, sqrtRatioX96, liquidity, tick, p)
		return pool
	}
)
