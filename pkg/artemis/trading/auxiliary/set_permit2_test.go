package artemis_trading_auxiliary

import "fmt"

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2() {
	nodeURL := ""
	ta := InitAuxiliaryTradingUtils(ctx, nodeURL, t.acc)
	t.Require().NotEmpty(ta)
	fmt.Println(ta.Account.PublicKey())
}
