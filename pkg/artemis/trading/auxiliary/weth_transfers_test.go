package artemis_trading_auxiliary

import (
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestWETH() {
	ta := InitAuxiliaryTradingUtils(ctx, t.goerliNode, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)
}
