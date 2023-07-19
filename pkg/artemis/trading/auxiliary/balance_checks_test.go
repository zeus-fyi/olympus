package artemis_trading_auxiliary

import (
	"fmt"

	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
)

func (t *ArtemisAuxillaryTestSuite) TestBalanceCheck() {
	ta := t.at1
	t.Require().NotEmpty(ta)

	bal, err := checkEthBalance(ctx, *ta.w3c())
	t.Require().Nil(err)
	t.Require().NotNil(bal)

	fmt.Println("bal", bal.String())

	bal, err = CheckAuxWETHBalance(ctx, *ta.w3c())
	t.Require().Nil(err)
	t.Require().NotNil(bal)
	fmt.Println("weth bal", bal.String())

	ok, err := checkEthBalanceGreaterThan(ctx, *ta.w3c(), artemis_eth_units.GweiMultiple(10000))
	t.Require().Nil(err)
	t.Require().True(ok)

	nonce, err := ta.getNonce(ctx)
	t.Require().NotNil(nonce)
	fmt.Println("nonce", nonce)

	ta = t.at2
	t.Require().NotEmpty(ta)

	bal, err = checkEthBalance(ctx, *ta.w3c())
	fmt.Println("bal", bal.String())

	bal, err = CheckAuxWETHBalance(ctx, *ta.w3c())
	t.Require().Nil(err)
	t.Require().NotNil(bal)
	fmt.Println("weth bal", bal.String())

	ok, err = checkEthBalanceGreaterThan(ctx, *ta.w3c(), artemis_eth_units.GweiMultiple(10000))
	t.Require().Nil(err)
	t.Require().True(ok)

	nonce, err = ta.getNonce(ctx)
	t.Require().NotNil(nonce)
	fmt.Println("nonce", nonce)
}
