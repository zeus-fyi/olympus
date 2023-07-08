package artemis_trading_auxiliary

import (
	"fmt"

	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

// TODO, test with goerli, then set mainnet, needs to approve first before weth transfer

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2() {
	nodeURL := t.goerliNode
	ta := InitAuxiliaryTradingUtils(ctx, nodeURL, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)
	fmt.Println(ta.Account.PublicKey())

	approveTx, err := ta.ApprovePermit2(ctx, artemis_trading_constants.WETH9ContractAddress)
	t.Require().Nil(err)
	t.Require().NotEmpty(approveTx)
	fmt.Println("approveTx", approveTx.Hash().String())
}
