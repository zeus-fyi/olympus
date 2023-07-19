package artemis_trading_auxiliary

import (
	"fmt"

	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
)

func (t *ArtemisAuxillaryTestSuite) TestMainnetBal() {
	age := encryption.NewAge(t.Tc.LocalAgePkey, t.Tc.LocalAgePubkey)
	t.acc3 = initTradingAccount2(ctx, age)
	w3aMainnet := web3_client.NewWeb3Client(t.mainnetNode, &t.acc3)
	w3aMainnet.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	atMainnet := InitAuxiliaryTradingUtils(ctx, w3aMainnet)
	token := getChainSpecificWETH(*atMainnet.w3c()).String()
	fmt.Println("token", token)
	t.Require().NotEmpty(w3aMainnet.Headers)

	t.Require().Equal("Bearer "+t.Tc.ProductionLocalTemporalBearerToken, w3aMainnet.Headers["Authorization"])
	t.Require().Equal(t.mainnetNode, atMainnet.nodeURL())
	t.Require().Equal(token, artemis_trading_constants.WETH9ContractAddress)

	bal, err := checkEthBalance(ctx, *atMainnet.w3c())
	t.Require().Nil(err)
	t.Require().NotNil(bal)

	fmt.Println("bal", bal.String())

	bal, err = CheckAuxWETHBalance(ctx, *atMainnet.w3c())
	t.Require().Nil(err)
	t.Require().NotNil(bal)
	fmt.Println("weth bal", bal.String())
}

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2Mainnet() {
	age := encryption.NewAge(t.Tc.LocalAgePkey, t.Tc.LocalAgePubkey)
	t.acc3 = initTradingAccount2(ctx, age)
	w3aMainnet := web3_client.NewWeb3Client(t.mainnetNode, &t.acc3)
	w3aMainnet.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	atMainnet := InitAuxiliaryTradingUtils(ctx, w3aMainnet)
	token := getChainSpecificWETH(*atMainnet.w3c()).String()
	fmt.Println("token", token)
	t.Require().NotEmpty(w3aMainnet.Headers)

	t.Require().Equal("Bearer "+t.Tc.ProductionLocalTemporalBearerToken, w3aMainnet.Headers["Authorization"])
	t.Require().Equal(t.mainnetNode, atMainnet.nodeURL())
	t.Require().Equal(token, artemis_trading_constants.WETH9ContractAddress)

	//approveTx, err := atMainnet.SetPermit2ApprovalForToken(ctx, token)
	//t.Require().Nil(err)
	//t.Require().NotEmpty(approveTx)
	//fmt.Println("approveTx", approveTx.Hash().String())
}

func (t *ArtemisAuxillaryTestSuite) TestFundAccount() {
	age := encryption.NewAge(t.Tc.LocalAgePkey, t.Tc.LocalAgePubkey)
	t.acc3 = initTradingAccount2(ctx, age)
	w3aMainnet := web3_client.NewWeb3Client(t.mainnetNode, &t.acc3)
	w3aMainnet.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	atMainnet := InitAuxiliaryTradingUtils(ctx, w3aMainnet)
	token := getChainSpecificWETH(*atMainnet.w3c()).String()
	fmt.Println("token", token)
	t.Require().NotEmpty(w3aMainnet.Headers)

	t.Require().Equal("Bearer "+t.Tc.ProductionLocalTemporalBearerToken, w3aMainnet.Headers["Authorization"])
	t.Require().Equal(t.mainnetNode, atMainnet.nodeURL())
	t.Require().Equal(token, artemis_trading_constants.WETH9ContractAddress)

	bal, err := checkEthBalance(ctx, *atMainnet.w3c())
	t.Require().Nil(err)
	t.Require().NotNil(bal)

	fmt.Println("bal", bal.String())

	bal, err = CheckAuxWETHBalance(ctx, *atMainnet.w3c())
	t.Require().Nil(err)
	t.Require().NotNil(bal)
	fmt.Println("weth bal", bal.String())

	// 0.4 eth
	toExchAmount := artemis_eth_units.GweiMultiple(400000000)

	cmd := t.testEthToWETH(&atMainnet, toExchAmount)
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
	ok, err := checkEthBalanceGreaterThan(ctx, *atMainnet.w3c(), toExchAmount)
	t.Require().Nil(err)
	t.Require().True(ok)

	tx, _, err := universalRouterCmdToTxBuilder(ctx, *atMainnet.w3c(), cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)

	//executedTx, err := atMainnet.universalRouterExecuteTx(ctx, tx)
	//t.Require().Nil(err)
	//fmt.Println("executedTx", executedTx.Hash().String())
}
