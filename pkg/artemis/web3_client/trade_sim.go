package web3_client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog/log"
)

func (u *UniswapClient) SimFullSandwichTrade(tf *TradeExecutionFlow) error {
	if u.DebugPrint {
		fmt.Println("executing full sandwich trade")
	}
	// this isn't included in trade gas costs since we amortize one time gas costs for permit2
	eb := EtherMultiple(10000)
	bal := (*hexutil.Big)(eb)
	err := u.Web3Client.SetBalance(ctx, u.Web3Client.PublicKey(), *bal)
	if err != nil {
		log.Err(err).Msg("error setting balance")
		return err
	}
	nv, _ := new(big.Int).SetString("0", 10)
	nvB := (*hexutil.Big)(nv)
	err = u.Web3Client.SetNonce(ctx, u.Web3Client.PublicKey(), *nvB)
	if err != nil {
		log.Err(err).Msg("error setting nonce")
		return err
	}
	max, _ := new(big.Int).SetString(MaxUINT, 10)
	approveTx, err := u.ApproveSpender(ctx, WETH9ContractAddress, Permit2SmartContractAddress, max)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}
	secondToken := tf.FrontRunTrade.AmountInAddr.String()
	if tf.FrontRunTrade.AmountInAddr.String() == WETH9ContractAddress {
		secondToken = tf.FrontRunTrade.AmountOutAddr.String()
	}
	approveTx, err = u.ApproveSpender(ctx, secondToken, Permit2SmartContractAddress, max)
	if err != nil {
		log.Warn().Interface("approveTx", approveTx).Err(err).Msg("error approving permit2")
		return err
	}

	err = u.Web3Client.MatchFrontRunTradeValues(tf)
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
	return u.MarkEndOfSimDueToErr(nil)
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
