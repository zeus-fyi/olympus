package artemis_trading_auxiliary

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestExecV2TradeCall() (*web3_client.UniversalRouterExecCmd, *types.Transaction) {
	ta := t.at1
	t.Require().Equal(t.goerliNode, ta.nodeURL())
	cmd, _ := t.testExecV2Trade(&ta, hestia_req_types.Goerli)
	tx, _, err := universalRouterCmdToTxBuilder(ctx, *ta.w3c(), cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().NotNil(cmd.Deadline)
	return cmd, tx
}

func (t *ArtemisAuxillaryTestSuite) TestExecV2TradeCallMainnetSim() (*web3_client.UniversalRouterExecCmd, *types.Transaction) {
	ta := t.simMainnetTrader
	err := ta.setupCleanSimEnvironment(ctx, 0)
	t.Require().Nil(err)
	cmd, _ := t.testExecV2Trade(&ta, hestia_req_types.Mainnet)
	tx, _, err := universalRouterCmdToTxBuilder(ctx, *ta.w3c(), cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().NotNil(cmd.Deadline)
	return cmd, tx
}

func (t *ArtemisAuxillaryTestSuite) TestExecV2TradeCallMainnetSimPepe() (*web3_client.UniversalRouterExecCmd, *artemis_eth_txs.Permit2Tx, *types.Transaction) {
	ta := t.simMainnetTrader
	err := ta.setupCleanSimEnvironment(ctx, 0)
	t.Require().Nil(err)
	cmd, pt := t.testExecV2TradePepe(&ta)
	tx, _, err := universalRouterCmdToTxBuilder(ctx, *ta.w3c(), cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().NotNil(cmd.Deadline)
	return cmd, pt, tx
}
func (t *ArtemisAuxillaryTestSuite) TestExecV2TradeCallMainnetSimPepeToWeth() (*web3_client.UniversalRouterExecCmd, *types.Transaction) {
	ta := t.simMainnetTrader
	err := ta.setupCleanSimEnvironment(ctx, 17681149)
	t.Require().Nil(err)

	//bullshitCoin := artemis_trading_constants.DaiContractAddress
	bullshitCoin := artemis_trading_constants.PepeContractAddr

	/*
		// 107298366909358958965695136227
		// 17681149
	*/
	bullshitAmount := artemis_eth_units.NewBigIntFromStr("107298366909358958965695136227")
	//bullshitAmount := artemis_eth_units.EtherMultiple(10)
	err = ta.w3c().SetERC20BalanceBruteForce(ctx, bullshitCoin, ta.w3c().PublicKey(), bullshitAmount)
	t.Require().Nil(err)

	_, err = ta.SetPermit2ApprovalForToken(ctx, bullshitCoin)
	t.Require().Nil(err)

	cmd, _ := t.testExecV2TradeFromPepeToWeth(&ta, bullshitCoin, bullshitAmount)
	tx, _, err := universalRouterCmdToTxBuilder(ctx, *ta.w3c(), cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().NotNil(cmd.Deadline)
	return cmd, tx
}

func (t *ArtemisAuxillaryTestSuite) testExecV2TradeFromPepeToWeth(ta *AuxiliaryTradingUtils, dumbCoin string, toExchAmount *big.Int) (*web3_client.UniversalRouterExecCmd, *artemis_eth_txs.Permit2Tx) {
	t.Require().NotEmpty(ta)
	wethAddr := getChainSpecificWETH(*ta.w3c())

	t.Require().Equal(ta.network(), hestia_req_types.Mainnet)
	t.Require().Equal(wethAddr, artemis_trading_constants.WETH9ContractAddressAccount)

	scamCoin := accounts.HexToAddress(dumbCoin)
	to := &artemis_trading_types.TradeOutcome{
		AmountIn:      toExchAmount,
		AmountInAddr:  scamCoin,
		AmountOutAddr: wethAddr,
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
	to.AmountOut = artemis_eth_units.NewBigIntFromStr("0")

	cmd, pt, err := GenerateTradeV2SwapFromTokenToToken(ctx, *ta.w3c(), nil, to)
	t.Require().Nil(err)
	t.Require().NotEmpty(pt)
	t.Require().NotEmpty(cmd)
	t.Require().Len(cmd.Commands, 2)
	for i, sc := range cmd.Commands {
		if i == 0 && sc.Command != artemis_trading_constants.Permit2Permit {
			t.Fail("expected Permit2Permit")
		}
		if i == 0 && sc.Command == artemis_trading_constants.Permit2Permit {
			t.Require().Equal(toExchAmount.String(), sc.DecodedInputs.(web3_client.Permit2PermitParams).Amount.String())
			t.Require().Equal(scamCoin.String(), sc.DecodedInputs.(web3_client.Permit2PermitParams).Token.String())
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
	return cmd, pt
}

func (t *ArtemisAuxillaryTestSuite) testExecV2Trade(ta *AuxiliaryTradingUtils, network string) (*web3_client.UniversalRouterExecCmd, *artemis_eth_txs.Permit2Tx) {
	t.Require().NotEmpty(ta)
	// owner account for permit2
	t.Require().Equal(ta.tradersAccount().Address().String(), ta.w3c().Address().String())
	t.Require().Equal(ta.tradersAccount().Address().String(), ta.w3a().Address().String())
	toExchAmount := artemis_eth_units.GweiMultiple(10000)
	wethAddr := getChainSpecificWETH(*ta.w3c())
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

	cmd, pt, err := GenerateTradeV2SwapFromTokenToToken(ctx, *ta.w3c(), nil, to)
	t.Require().Nil(err)
	t.Require().NotNil(pt)
	t.Require().Equal(ta.tradersAccount().PublicKey(), pt.Owner)
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
	return cmd, pt
}

// todo add permit2 nonce getter from db method
func (t *ArtemisAuxillaryTestSuite) testExecV2TradePepe(ta *AuxiliaryTradingUtils) (*web3_client.UniversalRouterExecCmd, *artemis_eth_txs.Permit2Tx) {
	t.Require().NotEmpty(ta)
	toExchAmount := artemis_eth_units.GweiMultiple(10000)
	wethAddr := getChainSpecificWETH(*ta.w3c())
	pepeAddr := artemis_trading_constants.PepeContractAddrAccount

	t.Require().Equal(ta.network(), hestia_req_types.Mainnet)
	t.Require().Equal(wethAddr, artemis_trading_constants.WETH9ContractAddressAccount)
	t.Require().Equal(pepeAddr, artemis_trading_constants.PepeContractAddrAccount)

	to := &artemis_trading_types.TradeOutcome{
		AmountIn:      toExchAmount,
		AmountInAddr:  wethAddr,
		AmountOutAddr: pepeAddr,
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

	cmd, pt, err := GenerateTradeV2SwapFromTokenToToken(ctx, *ta.w3c(), nil, to)
	t.Require().Nil(err)
	t.Require().NotEmpty(cmd)
	t.Require().Len(cmd.Commands, 2)
	t.Require().NotNil(pt)
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
	return cmd, pt
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
	mts := maxTradeSize()
	fmt.Println("oneEther     :", artemis_eth_units.Ether.String())
	fmt.Println("maxTradeSize :", mts.String())
	t.Assert().Equal(artemis_eth_units.Ether, artemis_eth_units.MulBigIntFromInt(mts, 4))
}
