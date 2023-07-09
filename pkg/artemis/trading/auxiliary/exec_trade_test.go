package artemis_trading_auxiliary

import (
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestExecV2Trade() {
	ta := InitAuxiliaryTradingUtils(ctx, t.goerliNode, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)
	toExchAmount := artemis_eth_units.GweiMultiple(1000)
	wethAddr := artemis_trading_constants.WETH9ContractAddressAccount
	daiAddr := artemis_trading_constants.DaiContractAddressAccount
	if ta.Network == hestia_req_types.Goerli {
		wethAddr = artemis_trading_constants.GoerliWETH9ContractAddressAccount
		daiAddr = artemis_trading_constants.GoerliDaiContractAddressAccount
	}
	to := &artemis_trading_types.TradeOutcome{
		AmountIn:      toExchAmount,
		AmountInAddr:  wethAddr,
		AmountOut:     artemis_eth_units.NewBigInt(0),
		AmountOutAddr: daiAddr,
	}
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
			t.Require().Equal(true, sc.DecodedInputs.(web3_client.V2SwapExactInParams).PayerIsSender)
			t.Require().Equal([]accounts.Address{to.AmountInAddr, to.AmountOutAddr}, sc.DecodedInputs.(web3_client.V2SwapExactInParams).Path)
			t.Require().NotEmpty(sc.DecodedInputs.(web3_client.V2SwapExactInParams).AmountOutMin)
			t.Require().Equal(artemis_trading_constants.UniversalRouterSenderAddress, sc.DecodedInputs.(web3_client.V2SwapExactInParams).To.String())
		}
	}

	tx, err := ta.universalRouterCmdBuilder(ctx, cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().NotNil(cmd.Deadline)

	//_, err = ta.universalRouterExecuteTx(ctx, tx)
	//t.Require().Nil(err)
	//fmt.Println("tx", tx.Hash().String())
}
