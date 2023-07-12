package artemis_trading_auxiliary

import (
	"fmt"
)

func (t *ArtemisAuxillaryTestSuite) TestGeneratePermit2Nonce() {
	//for i := 0; i < 10; i++ {
	//	val := ts.GeneratePermit2Nonce()
	//	fmt.Println(val)
	//}
}

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2Goerli() {
	t.Require().Equal(t.goerliNode, t.at2.nodeURL())
	t.testSetPermit2()
}

//func (t *ArtemisAuxillaryTestSuite) TestSetPermit2Mainnet() {
//	age := encryption.NewAge(t.Tc.LocalAgePkey, t.Tc.LocalAgePubkey)
//	t.acc3 = initTradingAccount2(ctx, age)
//	//	t.testSetPermit2(hestia_req_types.Mainnet, t.acc3)
//}

func (t *ArtemisAuxillaryTestSuite) testSetPermit2() {
	//t.Require().Equal(t.at1, acc)
	//t.Require().NotEmpty(t.at1)
	//fmt.Println(t.at1.getWeb3Client().PublicKey())
	token := t.at1.getChainSpecificWETH().String()
	fmt.Println("token", token)
	approveTx, err := t.at2.SetPermit2ApprovalForToken(ctx, t.at1.getChainSpecificWETH().String())
	t.Require().Nil(err)
	t.Require().NotEmpty(approveTx)
	fmt.Println("approveTx", approveTx.Hash().String())
}
