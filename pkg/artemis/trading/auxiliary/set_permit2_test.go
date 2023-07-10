package artemis_trading_auxiliary

import (
	"fmt"

	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestGeneratePermit2Nonce() {
	for i := 0; i < 10; i++ {
		val := ts.GeneratePermit2Nonce()
		fmt.Println(val)
	}
}

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2() {
	nodeURL := t.goerliNode
	ta := InitAuxiliaryTradingUtils(ctx, nodeURL, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)
	fmt.Println(ta.getWeb3Client().PublicKey())
	wethAddr := artemis_trading_constants.WETH9ContractAddress
	if ta.getWeb3Client().Network == hestia_req_types.Goerli {
		wethAddr = artemis_trading_constants.GoerliWETH9ContractAddress
	}
	approveTx, err := ta.SetPermit2ApprovalForToken(ctx, wethAddr)
	t.Require().Nil(err)
	t.Require().NotEmpty(approveTx)
	fmt.Println("approveTx", approveTx.Hash().String())
}
