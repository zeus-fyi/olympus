package artemis_trade_debugger

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	artemis_realtime_trading "github.com/zeus-fyi/olympus/pkg/artemis/trading"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_reporting "github.com/zeus-fyi/olympus/pkg/artemis/trading/reporting"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

const (
	TraderAccountSim = "0x000000641e80A183c8B736141cbE313E136bc8c6"
)

func (t *ArtemisTradeDebuggerTestSuite) TestActiveReplay() {
	bg, berr := artemis_reporting.GetBundlesProfitHistory(ctx, 0, 1)
	t.Require().Nil(berr)
	t.Require().NotNil(bg)
	w3c := web3_client.NewWeb3ClientFakeSigner("https://iris.zeus.fyi/v1beta/internal/")
	w3c.AddBearerToken(t.Tc.ProductionLocalBearerToken)
	w3c.AddSessionLockHeader(ZeusTestSessionLockHeaderValue)
	bgFull, berr := artemis_reporting.GetBundleSubmissionHistory(ctx, 0, 1)
	t.Require().Nil(berr)
	t.Require().NotNil(bgFull)

	for bundleHash, b := range bg.Map {
		if bundleHash != "0x6bd2a42c538bbfeb7df47c6125b199099bb34ec1be31ffd2852b2a210b25cbd9" {
			continue
		}
		bundleTx := b[0]
		txGroup := bgFull.Map[bundleHash]
		t.Require().Len(txGroup, 3)

		fr := bgFull.Map[bundleHash][0].EthTx
		t.Require().Equal(TraderAccountSim, fr.From)
		sandwich := bgFull.Map[bundleHash][1].EthTx
		t.Require().Equal(TraderAccountSim, sandwich.From)
		user := bgFull.Map[bundleHash][2].EthTx
		t.Require().NotEqual(TraderAccountSim, user.From)
		t.Require().Less(fr.Nonce, sandwich.Nonce)

		bundleTx.PrintBundleInfo()
		tf := bundleTx.TradeExecutionFlow
		//t.Require().Equal(int(tf.CurrentBlockNumber.Uint64()), bundleTx.EthTxReceipts.BlockNumber-1)

		err := w3c.ResetNetworkLocalToExtIrisTest(int(tf.CurrentBlockNumber.Uint64()))
		t.Require().Nil(err)
		t.Require().NotNil(tf)
		err = CheckExpectedReserves(context.Background(), w3c, tf)
		t.Require().Nil(err)

		prevSp := tf.SandwichPrediction
		tf.SandwichPrediction = web3_client.SandwichTradePrediction{}
		latestTf, err := BinarySearch(tf)
		t.Require().Nil(err)

		fmt.Println("prevSp.ExpectedProfit.Int64()", prevSp.ExpectedProfit.Int64(), "latestTf.SandwichPrediction.ExpectedProfit.Int64()", latestTf.SandwichPrediction.ExpectedProfit.Int64())
		//tf.FrontRunTrade.PrintDebug()
		//tf.SandwichTrade.PrintDebug()
		fmt.Println(latestTf.SandwichTrade.AmountOut.String())
		fmt.Println(latestTf.SandwichPrediction.ExpectedProfit.String())
		latestTf.CurrentBlockNumber = tf.CurrentBlockNumber
		latestTf.Tx = tf.Tx
		err = artemis_realtime_trading.ApplyMaxTransferTaxCore(ctx, latestTf)
		t.Require().Nil(err)
		fmt.Println("prevSp.ExpectedProfit.Int64()", prevSp.ExpectedProfit.Int64(), "latestTfPostTax.SandwichPrediction.ExpectedProfit.Int64()", latestTf.SandwichPrediction.ExpectedProfit.Int64())

		//tf.FrontRunTrade.PrintDebug()
		//tf.SandwichTrade.PrintDebug()

		/*
		      {"level":"info","txHash":"0x697fcd28179683530a0f509fbc97bb6affb3a4b2bcddd6245232f6eebafd6aa3","bn":17762536,"profitTokenAddress":"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		   	"sellAmount":79999999999999966,"tf.SandwichPrediction.ExpectedProfit":81310916184690182,"tf.SandwichTrade.AmountOut":"81310916184690182",
		   	"time":"2023-07-24T18:19:22-07:00","message":"ApplyMaxTransferTax: acceptable after tax"}

		*/
		accepted, err := artemis_trading_auxiliary.IsProfitTokenAcceptable(latestTf, nil)
		t.Require().Nil(err)
		t.Assert().True(accepted)

		fmt.Println("===============================================================================================================")
		fmt.Println("===============================================================================================================")
		tx, _, err := w3c.GetTxByHash(context.Background(), common.HexToHash(fr.TxHash))
		t.Require().Nil(err)
		_, decoded, err := web3_client.DecodeTxArgDataFromAbi(ctx, tx, artemis_oly_contract_abis.UniversalRouterNew)
		t.Require().Nil(err)
		ur, err := web3_client.NewDecodedUniversalRouterExecCmdFromMap(decoded, artemis_oly_contract_abis.UniversalRouterDecoder)
		t.Require().Nil(err)

		decodedCmd := ur.Commands[1].DecodedInputs.(web3_client.V2SwapExactInParams)
		fmt.Println("frontRun decoded amountIn", decodedCmd.AmountIn.String())
		t.Require().Equal(decodedCmd.AmountIn.String(), tf.FrontRunTrade.AmountIn.String())
		for _, addr := range decodedCmd.Path {
			fmt.Println(addr.String())
		}
		fmt.Println("payerIsSender", decodedCmd.PayerIsSender)
		fmt.Println("to", decodedCmd.To.String())

		wethBalStart, err := w3c.GetMainnetBalanceWETH(TraderAccountSim)
		t.Require().Nil(err)

		err = w3c.SendSignedTransaction(ctx, tx)
		t.Require().Nil(err)

		wethPostFrontRunBal, err := w3c.GetMainnetBalanceWETH(TraderAccountSim)
		t.Require().Nil(err)
		fmt.Println("wethBalStart", wethBalStart.String())
		fmt.Println("wethPostFrontRunBal", wethPostFrontRunBal.String())
		fmt.Println("wethPostFrontRunBal", artemis_eth_units.SubBigInt(wethPostFrontRunBal, wethBalStart).String())

		fmt.Println("===============================================================================================================")
		fmt.Println("===============================================================================================================")
		tx, _, err = w3c.GetTxByHash(context.Background(), common.HexToHash(user.TxHash))
		t.Require().Nil(err)

		//baseFee, err := w3c.GetBaseFee(context.Background())
		//t.Require().Nil(err)
		//fmt.Println("baseFee", baseFee.String())
		//
		//fmt.Println("tx.GasFeeCap().String())", tx.GasFeeCap().String())
		err = w3c.SendSignedTransaction(ctx, tx)
		t.Require().Nil(err)

		fmt.Println("===============================================================================================================")
		fmt.Println("===============================================================================================================")
		wethBalPreSandwich, err := w3c.GetMainnetBalanceWETH(TraderAccountSim)
		t.Require().Nil(err)

		otherTokenBalPreSandwich, err := w3c.ReadERC20TokenBalance(context.Background(), tf.SandwichTrade.AmountInAddr.String(), TraderAccountSim)
		t.Require().Nil(err)
		fmt.Println("otherTokenBalPreSandwich", otherTokenBalPreSandwich.String())

		expTrue := artemis_eth_units.IsXGreaterThanOrEqualToY(otherTokenBalPreSandwich, tf.SandwichTrade.AmountIn)
		t.Require().True(expTrue)

		tx, _, err = w3c.GetTxByHash(context.Background(), common.HexToHash(sandwich.TxHash))
		t.Require().Nil(err)
		fmt.Println("sandwich tx.Hash().String()", tx.Hash().String())
		fmt.Println("tf.SandwichTrade.AmountIn", tf.SandwichTrade.AmountIn.String())
		_, decoded, err = web3_client.DecodeTxArgDataFromAbi(ctx, tx, artemis_oly_contract_abis.UniversalRouterNew)
		t.Require().Nil(err)
		ur, err = web3_client.NewDecodedUniversalRouterExecCmdFromMap(decoded, artemis_oly_contract_abis.UniversalRouterDecoder)
		t.Require().Nil(err)
		t.Require().NotNil(ur)
		decodedCmd = ur.Commands[1].DecodedInputs.(web3_client.V2SwapExactInParams)
		fmt.Println("decodedBackRunAmountIn", decodedCmd.AmountIn.String())
		fmt.Println("decodedAmountOutMin", decodedCmd.AmountOutMin.String())
		fmt.Println("payerIsSender", decodedCmd.PayerIsSender)
		fmt.Println("to", decodedCmd.To.String())

		//t.Require().Equal(decodedCmd.To.String(), TraderAccountSim)

		for _, addr := range decodedCmd.Path {
			fmt.Println(addr.String())
		}
		err = w3c.SendSignedTransaction(ctx, tx)
		t.Require().Nil(err)

		wethPostSandwichBal, err := w3c.GetMainnetBalanceWETH(TraderAccountSim)
		t.Require().Nil(err)

		otherTokenBalPostSandwich, err := w3c.ReadERC20TokenBalance(context.Background(), tf.SandwichTrade.AmountInAddr.String(), TraderAccountSim)
		t.Require().Nil(err)
		fmt.Println("otherTokenBalPostSandwich", otherTokenBalPostSandwich.String())

		backRunWETHDiff := artemis_eth_units.SubBigInt(wethPostSandwichBal, wethBalPreSandwich)
		fmt.Println("wethDifferenceAfterSandwich", backRunWETHDiff.String())

		expTrue = artemis_eth_units.IsXGreaterThanOrEqualToY(backRunWETHDiff, decodedCmd.AmountOutMin)
		t.Require().True(expTrue)
		fmt.Println("wethBalPreSandwich", wethBalPreSandwich.String())
		fmt.Println("wethPostSandwichBal", wethPostSandwichBal.String())

		fmt.Println("===============================================================================================================")
		fmt.Println("===============================================================================================================")
		fmt.Println("txHash", tf.Tx.Hash().String())
		fmt.Println("wethDifferenceAfterSandwich", backRunWETHDiff.String())

		fmt.Println("prevSp.SandwichPrediction.ExpectedProfit.String()", prevSp.ExpectedProfit.String())
		fmt.Println("latestTf.SandwichPrediction.ExpectedProfit.String()", latestTf.SandwichPrediction.ExpectedProfit.String())

		fmt.Println("prevSp.SandwichTrade.AmountOut.String()", tf.SandwichTrade.AmountOut.String())
		fmt.Println("latestTf.SandwichTrade.AmountOut.String()", latestTf.SandwichTrade.AmountOut.String())

		fmt.Println("decodedAmountOutMin", decodedCmd.AmountOutMin.String())

		fmt.Println("actual bundleTx.Revenue Amount Out", bundleTx.Revenue)
		fmt.Println("tf.SandwichTrade.AmountInAddr.String()", tf.SandwichTrade.AmountInAddr.String())
		fmt.Println("===============================================================================================================")
		fmt.Println("===============================================================================================================")
	}
}
