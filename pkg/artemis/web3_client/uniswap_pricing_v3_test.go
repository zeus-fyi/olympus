package web3_client

import (
	"fmt"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/utils"
)

// example v3 pool: 0x4b5Ab61593A2401B1075b90c04cBCDD3F87CE011

func (s *Web3ClientTestSuite) TestUniswapV3DataFetcher() {
	tokenA := core_entities.NewToken(1, accounts.HexToAddress(WETH9ContractAddress), 18, "WETH", "Wrapped Ether")
	tokenB := core_entities.NewToken(1, accounts.HexToAddress(LooksTokenAddr), 18, "LOOKS", "LooksRare Token")

	p := UniswapPoolV3{
		PoolAddress: "0x4b5Ab61593A2401B1075b90c04cBCDD3F87CE011",
		Web3Actions: s.LocalHardhatMainnetUser.Web3Actions,
	}

	tick, err := p.GetTickFromContract(0)
	s.Require().NoError(err)
	s.Require().NotNil(tick)

	err = p.GetSlot0()
	s.Require().NoError(err)

	err = p.GetLiquidity()
	s.Require().NoError(err)

	v3Pool, err := entities.NewPool(tokenA, tokenB, constants.FeeMedium, p.Slot0.SqrtPriceX96, p.Liquidity, p.Slot0.Tick, &p)
	s.Require().NoError(err)
	s.Require().NotNil(v3Pool)

	inputAmount := core_entities.FromRawAmount(tokenA, big.NewInt(100))
	output, pool, err := v3Pool.GetOutputAmount(inputAmount, nil)
	s.Require().NoError(err)
	s.Require().NotNil(output)
	s.Require().NotNil(pool)

	fmt.Println(output.Numerator.String())
	fmt.Println(output.Currency.Name())
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

	/*
		func NewPool(tokenA, tokenB *entities.Token, fee constants.FeeAmount, sqrtRatioX96 *big.Int, liquidity *big.Int, tickCurrent int, ticks TickDataProvider) (*Pool, error) {

	*/
	// 	rpool_0_1, _    = NewPool(rtoken0, rtoken1, constants.FeeMedium, utils.EncodeSqrtRatioX96(constants.One, constants.One), big.NewInt(0), 0, nil)

	v3Pool, err := entities.NewPool(tokenA, tokenB, constants.FeeLow, nil, nil, 0, nil)
	s.Require().NoError(err)
	s.Require().NotNil(v3Pool)

	inputAmount := core_entities.FromRawAmount(tokenA, big.NewInt(100))
	output, pool, err := v3Pool.GetOutputAmount(inputAmount, nil)
	s.Require().NoError(err)
	s.Require().NotNil(output)
	s.Require().NotNil(pool)
}
