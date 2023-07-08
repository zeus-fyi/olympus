package artemis_trading_auxiliary

import (
	"fmt"

	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestBalanceCheck() {
	ta := InitAuxiliaryTradingUtils(ctx, t.goerliNode, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)

	bal, err := ta.checkAuxWETHBalance(ctx)
	t.Require().Nil(err)
	t.Require().NotNil(bal)
	fmt.Println("bal", bal.String())
}
