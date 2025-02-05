package artemis_trading_auxiliary

import (
	"context"

	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) testExecV2TradeFrontRun(ta *AuxiliaryTradingUtils) *web3_client.UniversalRouterExecCmd {
	t.Require().NotEmpty(ta)
	toExchAmount := artemis_eth_units.GweiMultiple(100)
	wethAddr := getChainSpecificWETH(*ta.w3c())
	daiAddr := artemis_trading_constants.DaiContractAddressAccount
	if ta.network() == hestia_req_types.Goerli {
		daiAddr = artemis_trading_constants.GoerliDaiContractAddressAccount
	}
	to := &artemis_trading_types.TradeOutcome{
		AmountIn:      toExchAmount,
		AmountInAddr:  wethAddr,
		AmountOutAddr: daiAddr,
	}
	path := []accounts.Address{to.AmountInAddr, to.AmountOutAddr}

	bn, berr := artemis_trading_cache.GetLatestBlockFromCacheOrProvidedSource(context.Background(), ta.U.Web3Client.Web3Actions)
	t.Require().Nil(berr)
	prices, err := artemis_uniswap_pricing.V2PairToPrices(ctx, bn, ta.U.Web3Client.Web3Actions, path)
	t.Require().Nil(err)
	amountOut, err := prices.GetQuoteUsingTokenAddr(to.AmountInAddr.String(), to.AmountIn)
	t.Require().Nil(err)
	t.Require().NotNil(amountOut)
	to.AmountOut = amountOut

	cmd, pt, err := GenerateTradeV2SwapFromTokenToToken(ctx, *ta.w3c(), nil, to)
	t.Require().Nil(err)
	t.Require().NotNil(pt)
	t.Require().NotEmpty(cmd)
	t.Require().Len(cmd.Commands, 2)
	for i, sc := range cmd.Commands {
		if i == 0 && sc.Command != artemis_trading_constants.Permit2Permit {
			t.Fail("expected Permit2Permit")
		}
		if i == 0 && sc.Command == artemis_trading_constants.Permit2Permit {
			t.Require().Equal(toExchAmount.String(), sc.DecodedInputs.(web3_client.Permit2PermitParams).Amount.String())
			t.Require().Equal(wethAddr.String(), sc.DecodedInputs.(web3_client.Permit2PermitParams).Token.String())
			t.Require().Equal(artemis_trading_constants.UniswapUniversalRouterAddressNew, sc.DecodedInputs.(web3_client.Permit2PermitParams).Spender.String())
		}
		if i == 1 && sc.Command != artemis_trading_constants.V2SwapExactIn {
			t.Fail("expected V2SwapExactIn")
		}
		if i == 0 && sc.Command == artemis_trading_constants.V2SwapExactIn {
			t.Require().Equal(toExchAmount.String(), sc.DecodedInputs.(web3_client.V2SwapExactInParams).AmountIn.String())
			t.Require().Equal(to.AmountOut.String(), sc.DecodedInputs.(web3_client.V2SwapExactInParams).AmountOutMin.String())
			t.Require().Equal(true, sc.DecodedInputs.(web3_client.V2SwapExactInParams).PayerIsSender)
			t.Require().Equal(path, sc.DecodedInputs.(web3_client.V2SwapExactInParams).Path)
			t.Require().NotEmpty(sc.DecodedInputs.(web3_client.V2SwapExactInParams).AmountOutMin)
			t.Require().Equal(artemis_trading_constants.UniversalRouterSenderAddress, sc.DecodedInputs.(web3_client.V2SwapExactInParams).To.String())
		}
	}
	return cmd
}
