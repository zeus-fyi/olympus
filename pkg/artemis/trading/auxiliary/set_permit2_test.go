package artemis_trading_auxiliary

import (
	"fmt"

	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestGeneratePermit2Nonce() {
	//for i := 0; i < 10; i++ {
	//	val := ts.GeneratePermit2Nonce()
	//	fmt.Println(val)
	//}
}

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2Goerli() {
	t.testSetPermit2(hestia_req_types.Goerli, t.acc2)
}

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2Mainnet() {
	age := encryption.NewAge(t.Tc.LocalAgePkey, t.Tc.LocalAgePubkey)
	t.acc3 = initTradingAccount2(ctx, age)

	//	t.testSetPermit2(hestia_req_types.Mainnet, t.acc3)
}

func (t *ArtemisAuxillaryTestSuite) testSetPermit2(network string, acc accounts.Account) {
	nodeURL := ""
	switch network {
	case hestia_req_types.Goerli:
		nodeURL = t.goerliNode
	case hestia_req_types.Mainnet:
		nodeURL = t.mainnetNode
	case hestia_req_types.Ephemery:
		t.Require().Fail("invalid network")
	default:
		t.Require().Fail("invalid network")
	}
	ta := InitAuxiliaryTradingUtils(ctx, nodeURL, network, acc)
	t.Require().NotEmpty(ta)
	fmt.Println(ta.getWeb3Client().PublicKey())
	approveTx, err := ta.SetPermit2ApprovalForToken(ctx, ta.getChainSpecificWETH().String())
	t.Require().Nil(err)
	t.Require().NotEmpty(approveTx)
	fmt.Println("approveTx", approveTx.Hash().String())
}
