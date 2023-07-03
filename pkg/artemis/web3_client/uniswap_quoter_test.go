package web3_client

import (
	"fmt"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	core_entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/utils"
)

func (s *Web3ClientTestSuite) TestUniswapQuoterV2() {
	factoryAddress := accounts.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984")
	tokenA := core_entities.NewToken(1, accounts.HexToAddress(WETH9ContractAddress), 18, "WETH", "Wrapped Ether")
	tokenB := core_entities.NewToken(1, accounts.HexToAddress(UsdCoinAddr), 6, "USDC", "USD Coin")
	result, err := utils.ComputePoolAddress(factoryAddress, tokenA, tokenB, constants.FeeMedium, "")
	s.Require().NoError(err)
	p := &uniswap_pricing.UniswapPoolV3{
		PoolAddress: result.String(),
		Web3Actions: s.LocalHardhatMainnetUser.Web3Actions,
	}
	err = p.GetSlot0(ctx)
	s.Require().NoError(err)

	err = p.GetLiquidity(ctx)
	s.Require().NoError(err)

	v3Pool, err := entities.NewPool(tokenA, tokenB, constants.FeeMedium, p.Slot0.SqrtPriceX96, p.Liquidity, p.Slot0.Tick, nil)
	s.Require().NoError(err)
	s.Require().NotNil(v3Pool)

	uni := InitUniswapClient(ctx, s.LocalHardhatMainnetUser)
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

	usdAmount := new(big.Int).Div(resp.AmountOut, new(big.Int).SetInt64(1000000))
	fmt.Println("usdcAmount", usdAmount.String())
}
