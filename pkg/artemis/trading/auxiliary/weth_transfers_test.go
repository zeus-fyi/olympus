package artemis_trading_auxiliary

import (
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestWETH() {
	ta := InitAuxiliaryTradingUtils(ctx, t.goerliNode, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)
	toExchAmount := artemis_eth_units.GweiMultiple(1000)
	cmd, err := ta.GenerateCmdToExchangeETHtoWETH(ctx, nil, toExchAmount, nil)
	t.Require().Nil(err)
	t.Require().NotEmpty(cmd)
	found := false
	for i, sc := range cmd.Commands {
		if i == 0 && sc.Command == artemis_trading_constants.WrapETH {
			found = true
			t.Require().NotNil(cmd.Payable.Amount)
			t.Require().Equal(toExchAmount.String(), cmd.Payable.Amount.String())
			t.Require().Equal(toExchAmount.String(), sc.DecodedInputs.(web3_client.WrapETHParams).AmountMin.String())
			t.Require().Equal(ta.Address().String(), sc.DecodedInputs.(web3_client.WrapETHParams).Recipient.String())
		}
	}
	t.Require().True(found)
	ok, err := ta.checkAuxEthBalanceGreaterThan(ctx, toExchAmount)
	t.Require().Nil(err)
	t.Require().True(ok)

	tx, err := ta.universalRouterCmdBuilder(ctx, cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)

	//_, err = ta.universalRouterExecuteTx(ctx, tx)
	//t.Require().Nil(err)
	//fmt.Println("tx", tx.Hash().String())
}

func (t *ArtemisAuxillaryTestSuite) TestUnwrapWETH() {
	ta := InitAuxiliaryTradingUtils(ctx, t.goerliNode, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)
	toExchAmount := artemis_eth_units.EtherMultiple(1)
	cmd, err := ta.GenerateCmdToExchangeWETHtoETH(ctx, nil, toExchAmount, nil)
	t.Require().Nil(err)
	t.Require().NotEmpty(cmd)
	found := false
	for i, sc := range cmd.Commands {
		if i == 0 && sc.Command == artemis_trading_constants.UnwrapWETH {
			found = true
			t.Require().Nil(cmd.Payable.Amount)
			t.Require().Equal(toExchAmount.String(), sc.DecodedInputs.(web3_client.UnwrapWETHParams).AmountMin.String())
			t.Require().Equal(ta.Address().String(), sc.DecodedInputs.(web3_client.UnwrapWETHParams).Recipient.String())
		}
	}
	t.Require().True(found)
	ok, err := ta.checkAuxWETHBalanceGreaterThan(ctx, toExchAmount)
	t.Require().Nil(err)
	t.Require().True(ok)

	tx, err := ta.universalRouterCmdBuilder(ctx, cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
}
