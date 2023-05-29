package web3_client

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (u *UniswapV2Client) SimFullSandwichTrade(tf *TradeExecutionFlowInBigInt) error {
	if u.DebugPrint {
		fmt.Println("executing full sandwich trade")
	}
	err := u.Web3Client.MatchFrontRunTradeValues(tf)
	if err != nil {
		log.Err(err).Msg("error executing front run balance setup")
		return err
	}
	_, err = u.ExecFrontRunTradeStepTokenTransfer(tf)
	if err != nil {
		log.Err(err).Msg("error executing front run trade step token transfer")
		return err
	}
	_, err = u.ExecUserTradeStep(tf)
	if err != nil {
		log.Err(err).Msg("error executing user trade step")
		return err
	}
	_, err = u.ExecSandwichTradeStepTokenTransfer(tf)
	if err != nil {
		log.Err(err).Msg("error executing sandwich trade step token transfer")
		return err
	}
	err = tf.GetAggregateGasUsage(ctx, u.Web3Client)
	if err != nil {
		log.Err(err).Msg("error getting aggregate gas usage")
		return err
	}
	return u.VerifyTradeResults(tf)
}

func (u *UniswapV2Client) SimUserOnlyTrade(tf *TradeExecutionFlowInBigInt) error {
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
