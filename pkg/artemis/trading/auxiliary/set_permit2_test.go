package artemis_trading_auxiliary

import (
	"fmt"

	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2() {
	nodeURL := t.goerliNode
	ta := InitAuxiliaryTradingUtils(ctx, nodeURL, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)
	fmt.Println(ta.Account.PublicKey())
	wethAddr := artemis_trading_constants.WETH9ContractAddress
	if ta.Network == hestia_req_types.Goerli {
		wethAddr = artemis_trading_constants.GoerliWETH9ContractAddress
	}
	approveTx, err := ta.ApprovePermit2(ctx, wethAddr)
	t.Require().Nil(err)
	t.Require().NotEmpty(approveTx)
	fmt.Println("approveTx", approveTx.Hash().String())
}

// don't think this is needed if we use permit2
func (t *ArtemisAuxillaryTestSuite) TestSetApproveUniversalRouter() {
	nodeURL := t.goerliNode
	ta := InitAuxiliaryTradingUtils(ctx, nodeURL, hestia_req_types.Goerli, t.acc)
	t.Require().NotEmpty(ta)
	fmt.Println(ta.Account.PublicKey())

	wethAddr := artemis_trading_constants.WETH9ContractAddress
	if ta.Network == hestia_req_types.Goerli {
		wethAddr = artemis_trading_constants.GoerliWETH9ContractAddress
	}
	approveTx, err := ta.ERC20ApproveSpender(ctx,
		wethAddr,
		artemis_trading_constants.UniswapUniversalRouterAddressNew,
		artemis_eth_units.MaxUINT)
	t.Require().Nil(err)
	t.Require().NotEmpty(approveTx)
	fmt.Println("approveTx", approveTx.Hash().String())
}
