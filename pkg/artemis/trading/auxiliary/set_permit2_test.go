package artemis_trading_auxiliary

import (
	"fmt"

	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
)

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2Goerli() {
	t.Require().Equal(t.goerliNode, t.at2.nodeURL())
	t.testSetPermit2()
}

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2Mainnet() {
	age := encryption.NewAge(t.Tc.LocalAgePkey, t.Tc.LocalAgePubkey)
	t.acc3 = initTradingAccount2(ctx, age)
	w3aMainnet := web3_client.NewWeb3Client(t.mainnetNode, &t.acc3)
	w3aMainnet.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	atMainnet := InitAuxiliaryTradingUtils(ctx, w3aMainnet)
	token := atMainnet.getChainSpecificWETH().String()
	fmt.Println("token", token)
	t.Require().NotEmpty(w3aMainnet.Headers)

	t.Require().Equal("Bearer "+t.Tc.ProductionLocalTemporalBearerToken, w3aMainnet.Headers["Authorization"])
	t.Require().Equal(t.mainnetNode, atMainnet.nodeURL())
	t.Require().Equal(token, artemis_trading_constants.WETH9ContractAddress)

	bal, err := atMainnet.checkEthBalance(ctx)
	t.Require().Nil(err)
	t.Require().NotNil(bal)

	fmt.Println("bal", bal.String())

	bal, err = atMainnet.CheckAuxWETHBalance(ctx)
	t.Require().Nil(err)
	t.Require().NotNil(bal)
	fmt.Println("weth bal", bal.String())

	//approveTx, err := atMainnet.SetPermit2ApprovalForToken(ctx, token)
	//t.Require().Nil(err)
	//t.Require().NotEmpty(approveTx)
	//fmt.Println("approveTx", approveTx.Hash().String())
}

func (t *ArtemisAuxillaryTestSuite) testSetPermit2() {
	//t.Require().Equal(t.at1, acc)
	//t.Require().NotEmpty(t.at1)
	//fmt.Println(t.at1.getWeb3Client().PublicKey())
	at := t.at1
	token := at.getChainSpecificWETH().String()
	fmt.Println("token", token)
	approveTx, err := at.SetPermit2ApprovalForToken(ctx, at.getChainSpecificWETH().String())
	t.Require().Nil(err)
	t.Require().NotEmpty(approveTx)
	fmt.Println("approveTx", approveTx.Hash().String())
}
