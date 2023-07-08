package artemis_trading_auxiliary

import (
	"fmt"

	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

// TODO, test with goerli, then set mainnet

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2() {
	nodeURL := t.goerliNode
	ta := InitAuxiliaryTradingUtils(ctx, nodeURL, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)
	fmt.Println(ta.Account.PublicKey())
}
