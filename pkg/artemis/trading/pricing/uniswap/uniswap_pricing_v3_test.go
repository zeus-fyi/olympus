package artemis_uniswap_pricing

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	artemis_multicall "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/multicall"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
)

var ctx = context.Background()

func (s *UniswapPricingTestSuite) TestMulticall3UniswapV3() {
	s.InitLocalConfigs()
	artemis_test_cache.InitLiveTestNetwork(s.Tc.QuikNodeURLS.TestRoute)
	wc := artemis_test_cache.LiveTestNetwork
	wc.Dial()
	defer wc.Close()
	p := UniswapV3Pair{
		Web3Actions: wc,
		PoolAddress: "0xF239009A101B6B930A527DEaaB6961b6E7deC8a6",
	}

	m3calls := []artemis_multicall.MultiCallElement{{
		Name: liquidity,
		Call: artemis_multicall.Call{
			Target:       common.HexToAddress(p.PoolAddress),
			AllowFailure: false,
			Data:         nil,
		},
		AbiFile:       v3PoolAbi,
		DecodedInputs: []interface{}{},
	}, {
		Name: slot0,
		Call: artemis_multicall.Call{
			Target:       common.HexToAddress(p.PoolAddress),
			AllowFailure: false,
			Data:         nil,
		},
		AbiFile:       v3PoolAbi,
		DecodedInputs: []interface{}{},
	}}
	m := artemis_multicall.Multicall3{
		Calls:   m3calls,
		Results: nil,
	}
	resp, err := m.PackAndCall(ctx, p.Web3Actions)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	for ind, respVal := range resp {
		fmt.Println(respVal.Success)
		fmt.Println(respVal.DecodedReturnData)
		switch ind {
		case 0:
			respSlice := respVal.DecodedReturnData
			for i, val := range respSlice {
				switch i {
				case 0:
					p.Liquidity = val.(*big.Int)
				}
			}
		case 1:
			respSlice := respVal.DecodedReturnData
			for i, val := range respSlice {
				switch i {
				case 0:
					p.Slot0.SqrtPriceX96 = val.(*big.Int)
				case 1:
					tmp := val.(*big.Int)
					p.Slot0.Tick = int(tmp.Int64())
				case 5:
					tmp := val.(uint8)
					p.Slot0.FeeProtocol = int(tmp)
				}
			}
		}
	}
	s.Require().NotEmpty(p.Liquidity)
	s.Require().NotEmpty(p.Slot0.SqrtPriceX96)
	s.Require().NotEmpty(p.Slot0.Tick)
	pc := UniswapV3Pair{
		Web3Actions: wc,
		PoolAddress: "0xF239009A101B6B930A527DEaaB6961b6E7deC8a6",
	}
	err = pc.GetLiquidityAndSlot0FromMulticall3(ctx)
	s.Require().NoError(err)
	s.Require().Equal(p.Liquidity, pc.Liquidity)
	s.Require().Equal(p.Slot0.SqrtPriceX96, pc.Slot0.SqrtPriceX96)
	s.Require().Equal(p.Slot0.Tick, pc.Slot0.Tick)
}
