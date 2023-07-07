package artemis_trading_auxiliary

import "fmt"

// TODO, test with goerli, then set mainnet

func (t *ArtemisAuxillaryTestSuite) TestSetPermit2() {
	nodeURL := ""
	ta := InitAuxiliaryTradingUtils(ctx, nodeURL, t.acc)
	t.Require().NotEmpty(ta)
	fmt.Println(ta.Account.PublicKey())
}
