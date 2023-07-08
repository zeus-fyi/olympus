package artemis_trading_auxiliary

import (
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestWETH() {
	ta := InitAuxiliaryTradingUtils(ctx, t.goerliNode, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)

	cmd, err := ta.GenerateCmdToExchangeETHtoWETH(ctx, nil, artemis_eth_units.EtherMultiple(1), nil)
	t.Require().Nil(err)
	t.Require().NotEmpty(cmd)
}
