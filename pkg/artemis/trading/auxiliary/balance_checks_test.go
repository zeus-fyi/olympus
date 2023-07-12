package artemis_trading_auxiliary

import (
	"fmt"
)

func (t *ArtemisAuxillaryTestSuite) TestBalanceCheck() {
	ta := t.at1
	t.Require().NotEmpty(ta)

	bal, err := ta.checkAuxWETHBalance(ctx)
	t.Require().Nil(err)
	t.Require().NotNil(bal)
	fmt.Println("bal", bal.String())

	ta = t.at2
	t.Require().NotEmpty(ta)

	bal, err = ta.checkAuxWETHBalance(ctx)
	t.Require().Nil(err)
	t.Require().NotNil(bal)
	fmt.Println("bal", bal.String())
}
