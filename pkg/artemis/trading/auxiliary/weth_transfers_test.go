package artemis_trading_auxiliary

import (
	"fmt"
	"math/big"

	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) testEthToWETH(ta *AuxiliaryTradingUtils, toExchAmount *big.Int) *web3_client.UniversalRouterExecCmd {
	t.Require().NotEmpty(ta)
	cmd, err := ta.GenerateCmdToExchangeETHtoWETH(ctx, nil, toExchAmount, nil)
	t.Require().Nil(err)
	t.Require().NotEmpty(cmd)
	return cmd
}

func (t *ArtemisAuxillaryTestSuite) TestWETH() {
	t.testWETH(hestia_req_types.Goerli)

	t.testWETH(hestia_req_types.Mainnet)
}

func (t *ArtemisAuxillaryTestSuite) testWETH(network string) {
	toExchAmount := artemis_eth_units.GweiMultiple(10000)

	ta := t.simMainnetTrader
	if network == hestia_req_types.Goerli {
		ta = t.at2
		t.Require().Equal(t.goerliNode, ta.nodeURL())
	} else {
		err := ta.setupCleanSimEnvironment(ctx, 0)
		t.Require().Nil(err)
	}

	cmd := t.testEthToWETH(&ta, toExchAmount)
	found := false
	for i, sc := range cmd.Commands {
		if i == 0 && sc.Command == artemis_trading_constants.WrapETH {
			found = true
			t.Require().NotNil(cmd.Payable.Amount)
			t.Require().Equal(artemis_trading_constants.UniswapUniversalRouterNewAddressAccount.String(), cmd.Payable.ToAddress.String())
			t.Require().Equal(toExchAmount.String(), cmd.Payable.Amount.String())
			t.Require().Equal(toExchAmount.String(), sc.DecodedInputs.(web3_client.WrapETHParams).AmountMin.String())
			t.Require().Equal(artemis_trading_constants.UniversalRouterSender, sc.DecodedInputs.(web3_client.WrapETHParams).Recipient.String())
		}
	}
	t.Require().True(found)
	ok, err := ta.checkEthBalanceGreaterThan(ctx, toExchAmount)
	t.Require().Nil(err)
	t.Require().True(ok)

	tx, _, err := ta.universalRouterCmdToTxBuilder(ctx, cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)

	if network == hestia_req_types.Goerli {
		_, err = ta.universalRouterExecuteTx(ctx, tx)
		t.Require().Nil(err)
		fmt.Println("tx", tx.Hash().String())
	}
}

func (t *ArtemisAuxillaryTestSuite) TestUnwrapWETH() {
	ta := t.at1
	t.Require().Equal(t.goerliNode, ta.nodeURL())
	t.Require().NotEmpty(ta)
	toExchAmount := artemis_eth_units.GweiMultiple(1000)
	cmd, err := ta.generateCmdToExchangeWETHtoETH(ctx, nil, toExchAmount, nil)
	t.Require().Nil(err)
	t.Require().NotEmpty(cmd)
	t.Require().Len(cmd.Commands, 2)
	t.Require().Nil(cmd.Payable.Amount)
	//wethAddr := artemis_trading_constants.WETH9ContractAddressAccount
	//if ta.Network == hestia_req_types.Goerli {
	//	wethAddr = artemis_trading_constants.GoerliWETH9ContractAddressAccount
	//}
	for i, sc := range cmd.Commands {
		//if i == 0 && sc.Command != artemis_trading_constants.Permit2Permit {
		//	t.Fail(fmt.Sprintf("expected %s, got %s", artemis_trading_constants.Permit2Permit, sc.Command))
		//}
		//if i == 0 {
		//	// token permissions
		//	t.Require().Equal(wethAddr.String(), sc.DecodedInputs.(web3_client.Permit2PermitParams).Token.String())
		//	t.Require().Equal(toExchAmount.String(), sc.DecodedInputs.(web3_client.Permit2PermitParams).Amount.String())
		//}
		if i == 1 && sc.Command != artemis_trading_constants.UnwrapWETH {
			t.Fail(fmt.Sprintf("expected %s, got %s", artemis_trading_constants.UnwrapWETH, sc.Command))
		}
		if i == 1 {
			t.Require().Equal(toExchAmount.String(), sc.DecodedInputs.(web3_client.UnwrapWETHParams).AmountMin.String())
			t.Require().Equal(artemis_trading_constants.UniversalRouterSender, sc.DecodedInputs.(web3_client.UnwrapWETHParams).Recipient.String())
		}
	}
	ok, err := CheckAuxWETHBalanceGreaterThan(ctx, *ta.w3c(), toExchAmount)
	t.Require().Nil(err)
	t.Require().True(ok)

	tx, _, err := ta.universalRouterCmdToTxBuilder(ctx, cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)

	//_, err = ta.universalRouterExecuteTx(ctx, tx)
	//t.Require().Nil(err)
	//fmt.Println("tx", tx.Hash().String())
}

/*
	t.Require().Equal(wethAddr.String(), sc.DecodedInputs.(web3_client.Permit2TransferFromParams).TokenPermissions.Token.String())
		t.Require().Equal(toExchAmount.String(), sc.DecodedInputs.(web3_client.Permit2TransferFromParams).TokenPermissions.Amount.String())
		// transfer details
		t.Require().Equal(toExchAmount.String(), sc.DecodedInputs.(web3_client.Permit2TransferFromParams).Permit2SignatureTransferDetails.RequestedAmount.String())
		t.Require().Equal(artemis_trading_constants.UniswapUniversalRouterAddressNew, sc.DecodedInputs.(web3_client.Permit2TransferFromParams).Permit2SignatureTransferDetails.To.String())
		// owner details
		t.Require().Equal(ta.Address().String(), sc.DecodedInputs.(web3_client.Permit2TransferFromParams).Owner.String())
		t.Require().NotNil(sc.DecodedInputs.(web3_client.Permit2TransferFromParams).Signature)
*/
