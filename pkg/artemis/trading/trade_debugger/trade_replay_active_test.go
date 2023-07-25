package artemis_trade_debugger

import (
	"context"

	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	artemis_reporting "github.com/zeus-fyi/olympus/pkg/artemis/trading/reporting"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

//w3cArchive := web3_client.NewWeb3ClientFakeSigner("https://eth-mainnet.g.alchemy.com/v2/cdVqiD1oZGvBiNEU8rDYt5kb6Q24nBMB")

func (t *ArtemisTradeDebuggerTestSuite) TestActiveReplay() {
	bg, berr := artemis_reporting.GetBundlesProfitHistory(ctx, 0, 1)
	t.Assert().Nil(berr)
	t.Assert().NotNil(bg)
	w3c := web3_client.NewWeb3ClientFakeSigner("https://eth.zeus.fyi")
	w3c.AddBearerToken(t.Tc.ProductionLocalBearerToken)

	w3cArchive := web3_client.NewWeb3ClientFakeSigner("https://eth-mainnet.g.alchemy.com/v2/cdVqiD1oZGvBiNEU8rDYt5kb6Q24nBMB")

	for _, b := range bg.Map {
		for _, bundleTx := range b {
			bundleTx.PrintBundleInfo()
			err := w3c.ResetNetworkLocalToExtIrisTest(bundleTx.EthTxReceipts.BlockNumber - 1)
			t.Require().Nil(err)
			tf := bundleTx.TradeExecutionFlow
			t.Require().NotNil(tf)
			err = CheckExpectedReserves(context.Background(), w3cArchive, bundleTx.TradeExecutionFlow)
			t.Require().Nil(err)

			_, err = BinarySearch(*tf)
			t.Require().Nil(err)
			//tf.FrontRunTrade.PrintDebug()
			//tf.SandwichTrade.PrintDebug()

			err = artemis_realtime_trading.ApplyMaxTransferTaxCore(ctx, tf)
			t.Require().Nil(err)
			//tf.FrontRunTrade.PrintDebug()
			//tf.SandwichTrade.PrintDebug()
			/*
			      {"level":"info","txHash":"0x697fcd28179683530a0f509fbc97bb6affb3a4b2bcddd6245232f6eebafd6aa3","bn":17762536,"profitTokenAddress":"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
			   	"sellAmount":79999999999999966,"tf.SandwichPrediction.ExpectedProfit":81310916184690182,"tf.SandwichTrade.AmountOut":"81310916184690182",
			   	"time":"2023-07-24T18:19:22-07:00","message":"ApplyMaxTransferTax: acceptable after tax"}

			*/
			accepted, err := artemis_trading_auxiliary.IsProfitTokenAcceptable(tf, nil)
			t.Require().Nil(err)
			t.Assert().True(accepted)
		}
	}
}
