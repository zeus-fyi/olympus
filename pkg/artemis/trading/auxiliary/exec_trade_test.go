package artemis_trading_auxiliary

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestExecV2TradeCall() (*web3_client.UniversalRouterExecCmd, *types.Transaction) {
	ta := t.at2
	t.Require().Equal(t.goerliNode, t.at2.nodeURL())
	cmd := t.testExecV2Trade(&ta, hestia_req_types.Goerli)
	tx, err := ta.universalRouterCmdToTxBuilder(ctx, cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().NotNil(cmd.Deadline)
	return cmd, tx
}

func (t *ArtemisAuxillaryTestSuite) TestExecV2TradeCallMainnetSim() (*web3_client.UniversalRouterExecCmd, *types.Transaction) {
	ta := t.simMainnetTrader
	err := ta.setupCleanSimEnvironment(ctx)
	t.Require().Nil(err)
	cmd := t.testExecV2Trade(&ta, hestia_req_types.Mainnet)
	tx, err := ta.universalRouterCmdToTxBuilder(ctx, cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().NotNil(cmd.Deadline)
	return cmd, tx
}

// todo add permit2 nonce getter from db method
func (t *ArtemisAuxillaryTestSuite) testExecV2Trade(ta *AuxiliaryTradingUtils, network string) *web3_client.UniversalRouterExecCmd {
	t.Require().NotEmpty(ta)
	toExchAmount := artemis_eth_units.GweiMultiple(10000)
	wethAddr := ta.getChainSpecificWETH()
	daiAddr := artemis_trading_constants.DaiContractAddressAccount
	if ta.network() == hestia_req_types.Goerli {
		daiAddr = artemis_trading_constants.GoerliDaiContractAddressAccount
	}
	if network == hestia_req_types.Goerli {
		t.Require().Equal(t.goerliNode, ta.w3c().NodeURL)
		t.Require().Equal(ta.network(), hestia_req_types.Goerli)
		t.Require().Equal(wethAddr, artemis_trading_constants.GoerliWETH9ContractAddressAccount)
		t.Require().Equal(daiAddr, artemis_trading_constants.GoerliDaiContractAddressAccount)
	} else {
		t.Require().Equal(ta.network(), hestia_req_types.Mainnet)
		t.Require().Equal(wethAddr, artemis_trading_constants.WETH9ContractAddressAccount)
		t.Require().Equal(daiAddr, artemis_trading_constants.DaiContractAddressAccount)
	}

	to := &artemis_trading_types.TradeOutcome{
		AmountIn:      toExchAmount,
		AmountInAddr:  wethAddr,
		AmountOutAddr: daiAddr,
	}
	path := []accounts.Address{to.AmountInAddr, to.AmountOutAddr}
	prices, err := artemis_uniswap_pricing.V2PairToPrices(ctx, *ta.w3a(), path)
	t.Require().Nil(err)
	t.Require().NotEmpty(prices)
	fmt.Println("testExecV2Trade: prices", prices.Reserve0.String(), prices.Reserve1.String())
	amountOut, err := prices.GetQuoteUsingTokenAddr(to.AmountInAddr.String(), to.AmountIn)
	t.Require().Nil(err)
	t.Require().NotNil(amountOut)
	fmt.Println("testExecV2Trade: amountOut", amountOut.String())
	to.AmountOut = amountOut
	cmd, err := ta.GenerateTradeV2SwapFromTokenToToken(ctx, nil, to)
	t.Require().Nil(err)
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

//func (t *ArtemisAuxillaryTestSuite) TestExecV2TradeExec() {
//	ta, _, tx := t.TestExecV2TradeCall()
//	_, err := ta.universalRouterExecuteTx(ctx, tx)
//	t.Require().Nil(err)
//	fmt.Println("tx", tx.Hash().String())
//}

func (t *ArtemisAuxillaryTestSuite) TestMaxTradeSize() {
	ta := t.at1
	t.Require().Equal(t.goerliNode, t.at1.nodeURL())
	t.Require().NotEmpty(ta)
	mts := ta.maxTradeSize()
	fmt.Println("oneEther     :", artemis_eth_units.Ether.String())
	fmt.Println("maxTradeSize :", mts.String())
	t.Assert().Equal(artemis_eth_units.Ether, artemis_eth_units.MulBigIntFromInt(mts, 4))
}
