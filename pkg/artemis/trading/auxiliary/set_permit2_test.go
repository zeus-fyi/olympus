package artemis_trading_auxiliary

import (
	"fmt"
)

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2Goerli() {
	t.Require().Equal(t.goerliNode, t.at2.nodeURL())
	t.testSetPermit2()
}

func (t *ArtemisAuxillaryTestSuite) testSetPermit2() {
	//t.Require().Equal(t.at1, acc)
	//t.Require().NotEmpty(t.at1)
	//fmt.Println(t.at1.getWeb3Client().PublicKey())
	at := t.at1
	token := getChainSpecificWETH(*at.w3c()).String()
	fmt.Println("token", token)
	approveTx, err := at.SetPermit2ApprovalForToken(ctx, getChainSpecificWETH(*at.w3c()).String())
	t.Require().Nil(err)
	t.Require().NotEmpty(approveTx)
	fmt.Println("approveTx", approveTx.Hash().String())
}
