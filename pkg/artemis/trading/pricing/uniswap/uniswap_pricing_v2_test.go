package artemis_uniswap_pricing

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis/v9"
	rdb "github.com/zeus-fyi/olympus/datastores/redis/apps"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_multicall "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/multicall"
	artemis_utils "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/utils"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
)

func (s *UniswapPricingTestSuite) TestRedisCache() {
	s.InitLocalConfigs()
	artemis_test_cache.InitLiveTestNetwork(s.Tc.QuikNodeURLS.TestRoute)
	wc := artemis_test_cache.LiveTestNetwork
	wc.Dial()
	defer wc.Close()
	redisOpts := redis.Options{
		Network: "",
		Addr:    "localhost:6379",
	}
	redisCache.Client = rdb.InitRedis(ctx, redisOpts)
	artemis_trading_cache.ReadRedis.Client = rdb.InitRedis(ctx, redisOpts)
	artemis_trading_cache.WriteRedis.Client = rdb.InitRedis(ctx, redisOpts)
	err := redisCache.AddV2PairToNextLookupSet(ctx, 0, "0x6C0207FB939596eCC63b4549ce22dFFF4c928216", "")
	s.Require().Nil(err)

	err = redisCache.AddV2PairToNextLookupSet(ctx, 0, "0xDE2FCae812b9EDda8d658bBBAa60ABB972B4D468", "")
	s.Require().Nil(err)

	addresses, err := redisCache.GetV2PairsToMulticall(ctx, 1)
	s.Require().Nil(err)
	s.Require().NotEmpty(addresses)
	s.Require().Equal(2, len(addresses))
}

func (s *UniswapPricingTestSuite) TestMulticall3UniswapV2Batch() {
	s.InitLocalConfigs()
	artemis_test_cache.InitLiveTestNetwork(s.Tc.QuikNodeURLS.TestRoute)
	wc := artemis_test_cache.LiveTestNetwork
	wc.Dial()
	defer wc.Close()
	p0 := "0x6C0207FB939596eCC63b4549ce22dFFF4c928216"
	p1 := "0xDE2FCae812b9EDda8d658bBBAa60ABB972B4D468"

	m3calls := []artemis_multicall.MultiCallElement{{
		Name: getReserves,
		Call: artemis_multicall.Call{
			Target:       common.HexToAddress(p0),
			AllowFailure: false,
			Data:         nil,
		},
		AbiFile:       v2ABI,
		DecodedInputs: []interface{}{},
	}, {
		Name: getReserves,
		Call: artemis_multicall.Call{
			Target:       common.HexToAddress(p1),
			AllowFailure: false,
			Data:         nil,
		},
		AbiFile:       v2ABI,
		DecodedInputs: []interface{}{},
	}}
	m := artemis_multicall.Multicall3{
		Calls:   m3calls,
		Results: nil,
	}
	resp, err := m.PackAndCall(ctx, wc)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(2, len(resp))

	pairOne := &UniswapV2Pair{}
	pairTwo := &UniswapV2Pair{}
	pairs := []*UniswapV2Pair{pairOne, pairTwo}
	for i, respVal := range resp {
		respSlice := respVal.DecodedReturnData
		p := &UniswapV2Pair{}
		reserve0, rerr := artemis_utils.ParseBigInt(respSlice[0])
		s.Require().NoError(rerr)
		p.Reserve0 = reserve0
		reserve1, rerr := artemis_utils.ParseBigInt(respSlice[1])
		s.Require().NoError(rerr)
		p.Reserve1 = reserve1

		fmt.Println("reserve0", reserve0.String())
		fmt.Println("reserve1", reserve1.String())
		blockTimestampLast, rerr := artemis_utils.ParseBigInt(respSlice[2])
		s.Require().NoError(rerr)
		p.BlockTimestampLast = blockTimestampLast
		pairs[i] = p
	}

	s.Require().NotEmpty(pairs[0])
	s.Require().NotEmpty(pairs[1])

	mcallResp, err := GetBatchPairContractPricesViaMulticall3(ctx, wc, p0, p1)
	s.Require().NoError(err)
	s.Require().NotNil(mcallResp)
	s.Require().Equal(2, len(mcallResp))

	for i, respVal := range mcallResp {
		s.Require().Equal(respVal.Reserve0.String(), pairs[i].Reserve0.String())
		s.Require().Equal(respVal.Reserve1.String(), pairs[i].Reserve1.String())
	}
}

func (s *UniswapPricingTestSuite) TestPricingImpact() {
	reserve0, _ := new(big.Int).SetString("400000", 10)  // TokenB
	reserve1, _ := new(big.Int).SetString("1200000", 10) // TokenA
	token0Addr, token1Addr := artemis_utils.StringsToAddresses(artemis_trading_constants.PepeContractAddr, WETH9ContractAddress)
	mockPairResp := UniswapV2Pair{
		KLast:    big.NewInt(0),
		Token0:   token0Addr,
		Token1:   token1Addr,
		Reserve0: reserve0,
		Reserve1: reserve1,
	}
	to, reservesToken0, reservesToken1 := mockPairResp.PriceImpactToken1BuyToken0(big.NewInt(3000))
	fmt.Println("to.AmountOut", to.AmountOut.String())
	fmt.Println("reservesToken0", reservesToken0.String())
	fmt.Println("reservesToken1", reservesToken1.String())
	s.Assert().Equal(big.NewInt(399006), reservesToken0)
	s.Assert().Equal(big.NewInt(1203000), reservesToken1)

	reserve0, _ = new(big.Int).SetString("400000", 10)  // TokenB
	reserve1, _ = new(big.Int).SetString("1200000", 10) // TokenA
	mockPairResp = UniswapV2Pair{
		KLast:    big.NewInt(0),
		Token0:   token0Addr,
		Token1:   token1Addr,
		Reserve0: reserve0,
		Reserve1: reserve1,
	}
	to, reservesToken0, reservesToken1 = mockPairResp.PriceImpactToken0BuyToken1(big.NewInt(1000))
	fmt.Println("to.AmountOut", to.AmountOut.String())
	fmt.Println("reservesToken0", reservesToken0.String())
	fmt.Println("reservesToken1", reservesToken1.String())
}

func (s *UniswapPricingTestSuite) TestUniswapSortTokens() {
	p := UniswapV2Pair{}
	err := p.PairForV2("0x8647Ae4E646cd3CE37FdEB4591b0A7928254bb73", WETH9ContractAddress)
	s.Require().Nil(err)
	s.Require().Equal("0xa20d8188e6844e41C7A8d7222D2A38611aa84784", p.PairContractAddr)
	s.Require().Equal(p.Token0.String(), "0x8647Ae4E646cd3CE37FdEB4591b0A7928254bb73")
	s.Require().Equal(p.Token1.String(), WETH9ContractAddress)
}
