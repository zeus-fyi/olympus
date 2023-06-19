package web3_client

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

func (u *UniswapClient) SimFullSandwichTrade(tf *TradeExecutionFlow) error {
	if u.DebugPrint {
		fmt.Println("executing full sandwich trade")
	}
	u.TradeAnalysisReport.TradeMethod = tf.Trade.TradeMethod
	err := u.Web3Client.MatchFrontRunTradeValues(tf)
	if err != nil {
		u.TradeFailureReport.EndStage = "executing front run balance setup"
		log.Err(err).Msg("error executing front run balance setup")
		return u.MarkEndOfSimDueToErr(err)
	}
	_, err = u.ExecFrontRunTradeStepTokenTransfer(tf)
	if err != nil {
		u.TradeFailureReport.EndStage = "executing front run trade"
		log.Err(err).Msg("error executing front run trade step token transfer")
		return u.MarkEndOfSimDueToErr(err)
	}
	_, err = u.ExecUserTradeStep(tf)
	if err != nil {
		u.TradeFailureReport.EndStage = "executing user trade step"
		log.Err(err).Msg("error executing user trade step")
		return u.MarkEndOfSimDueToErr(err)
	}
	_, err = u.ExecSandwichTradeStepTokenTransfer(tf)
	if err != nil {
		u.TradeFailureReport.EndStage = "executing sandwich trade"
		log.Err(err).Msg("error executing sandwich trade step token transfer")
		return u.MarkEndOfSimDueToErr(err)
	}
	err = tf.GetAggregateGasUsage(ctx, u.Web3Client)
	if err != nil {
		u.TradeFailureReport.EndStage = "post trade getting gas usage"
		log.Err(err).Msg("error getting aggregate gas usage")
		return u.MarkEndOfSimDueToErr(err)
	}
	err = u.VerifyTradeResults(tf)
	if err != nil {
		u.TradeFailureReport.EndStage = "verifying trade results"
		log.Err(err).Msg("error verifying trade results")
		return u.MarkEndOfSimDueToErr(err)
	}
	if !u.TestMode {
		return u.MarkEndOfSimDueToErr(nil)
	}
	return nil
}

func (u *UniswapClient) SimFrontRunTradeOnly(tf *TradeExecutionFlow) error {
	err := u.Web3Client.MatchFrontRunTradeValues(tf)
	if err != nil {
		log.Err(err).Msg("error executing front run balance setup")
		return err
	}
	fmt.Println("amountIn", tf.FrontRunTrade.AmountIn.String())
	tokenBal, err := u.Web3Client.ReadERC20TokenBalance(context.Background(), tf.FrontRunTrade.AmountInAddr.String(), u.Web3Client.PublicKey())
	if err != nil {
		log.Err(err).Msg("error reading token balance")
		return err
	}
	tokenBalOut, err := u.Web3Client.ReadERC20TokenBalance(context.Background(), tf.FrontRunTrade.AmountOutAddr.String(), u.Web3Client.PublicKey())
	if err != nil {
		log.Err(err).Msg("error reading token balance")
		return err
	}
	fmt.Println("token balance amountIn", tokenBal.String())
	fmt.Println("token balance amountOut", tokenBalOut.String())
	_, err = u.ExecFrontRunTradeStepTokenTransfer(tf)
	if err != nil {
		log.Err(err).Msg("error executing front run trade step token transfer")
		return err
	}
	fmt.Println("post-trade")

	tokenBal, err = u.Web3Client.ReadERC20TokenBalance(context.Background(), tf.FrontRunTrade.AmountInAddr.String(), u.Web3Client.PublicKey())
	if err != nil {
		log.Err(err).Msg("error reading token balance")
		return err
	}
	tokenBalOut, err = u.Web3Client.ReadERC20TokenBalance(context.Background(), tf.FrontRunTrade.AmountOutAddr.String(), u.Web3Client.PublicKey())
	if err != nil {
		log.Err(err).Msg("error reading token balance")
		return err
	}
	fmt.Println("token balance amountIn", tokenBal.String())
	fmt.Println("token balance amountOut", tokenBalOut.String())
	return err
}

func (u *UniswapClient) SimUserOnlyTrade(tf *TradeExecutionFlow) error {
	if u.DebugPrint {
		fmt.Println("executing stand alone user trade")
	}

	to, err := tf.InitialPair.PriceImpact(tf.UserTrade.AmountInAddr, tf.UserTrade.AmountIn)
	if err != nil {
		log.Err(err).Msg("error executing user trade prediction")
		return err
	}
	tf.UserTrade.AmountOut = to.AmountOut
	_, err = u.ExecUserTradeStep(tf)
	if err != nil {
		log.Err(err).Msg("error executing user trade step")
		return err
	}

	fmt.Println("amountInAddr", to.AmountInAddr.String())
	fmt.Println("amountIn", to.AmountIn.String())
	fmt.Println("amountOutAddr", to.AmountOutAddr.String())
	fmt.Println("amountOut", to.AmountOut.String())
	return nil
}
