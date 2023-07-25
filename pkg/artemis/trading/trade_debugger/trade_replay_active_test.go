package artemis_trade_debugger

import (
	"context"

	artemis_reporting "github.com/zeus-fyi/olympus/pkg/artemis/trading/reporting"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

//w3cArchive := web3_client.NewWeb3ClientFakeSigner("https://eth-mainnet.g.alchemy.com/v2/cdVqiD1oZGvBiNEU8rDYt5kb6Q24nBMB")

func (t *ArtemisTradeDebuggerTestSuite) TestActiveReplay() {
	bg, err := artemis_reporting.GetBundlesProfitHistory(ctx, 0, 1)
	t.Assert().Nil(err)
	t.Assert().NotNil(bg)
	w3c := web3_client.NewWeb3ClientFakeSigner("https://eth.zeus.fyi")
	w3c.AddBearerToken(t.Tc.ProductionLocalBearerToken)

	w3cArchive := web3_client.NewWeb3ClientFakeSigner("https://eth-mainnet.g.alchemy.com/v2/cdVqiD1oZGvBiNEU8rDYt5kb6Q24nBMB")

	for _, b := range bg.Map {
		for _, bundleTx := range b {
			bundleTx.PrintBundleInfo()
			err = w3c.ResetNetworkLocalToExtIrisTest(bundleTx.EthTxReceipts.BlockNumber - 1)
			t.Require().Nil(err)
			err = CheckExpectedReserves(context.Background(), w3cArchive, bundleTx.TradeExecutionFlow)
			t.Require().Nil(err)
		}
	}
	// TODO, get raw tx
	// reset block state
	// re-calc binary search
}
