package web3_client

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

func (u *UniswapClient) ExecSwap(pair UniswapV2Pair, to *TradeOutcome) (*web3_actions.SendContractTxPayload, error) {
	scInfo, _, err := LoadSwapAbiPayload(pair.PairContractAddr)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	tokenNum := pair.GetTokenNumber(to.AmountInAddr)
	if tokenNum == 0 {
		scInfo.Params = []interface{}{"0", to.AmountOut, u.Web3Client.Address(), []byte{}}
	} else {
		scInfo.Params = []interface{}{to.AmountOut, "0", u.Web3Client.Address(), []byte{}}
	}
	// TODO implement better gas estimation
	scInfo.GasLimit = 3000000
	signedTx, err := u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	to.AddTxHash(accounts.Hash(signedTx.Hash()))
	err = u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return &web3_actions.SendContractTxPayload{}, err
	}
	return &scInfo, nil
}

func (u *UniswapClient) ExecFrontRunTrade(tf TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(tf.InitialPair, &tf.FrontRunTrade)
}

func (u *UniswapClient) ExecSandwichTrade(tf TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(tf.InitialPair, &tf.SandwichTrade)
}

func (u *UniswapClient) ExecFrontRunTradeStepTokenTransfer(tf *TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	if u.DebugPrint {
		fmt.Println("executing front run token transfer")
	}
	_, _ = u.FrontRunTradeGetAmountsOut(tf)
	ethBal, err := u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting pre trade eth balance")
		return nil, err
	}
	tf.FrontRunTrade.PreTradeEthBalance = ethBal
	err = u.RouterApproveAndSend(ctx, &tf.FrontRunTrade, tf.InitialPair.PairContractAddr)
	if err != nil {
		return nil, err
	}
	ethBal, err = u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting post trade eth balance")
		return nil, err
	}
	tf.FrontRunTrade.PostTradeEthBalance = ethBal
	return u.ExecFrontRunTradeStep(tf)
}

func (u *UniswapClient) ExecUserTradeStep(tf *TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	if u.DebugPrint {
		fmt.Println("executing user trade")
	}
	_, _ = u.UserTradeGetAmountsOut(tf)
	sender := types.LatestSignerForChainID(tf.Tx.ChainId())
	from, err := sender.Sender(tf.Tx)
	if err != nil {
		log.Err(err).Msg("error getting sender")
		return nil, err
	}
	ethBal, err := u.Web3Client.GetBalance(ctx, from.String(), nil)
	if err != nil {
		log.Err(err).Msg("error getting pre trade eth balance")
		return nil, err
	}
	tf.UserTrade.PreTradeEthBalance = ethBal
	scInfo, err := u.ExecTradeByMethod(tf)
	if err != nil {
		return nil, err
	}
	ethBal, err = u.Web3Client.GetBalance(ctx, from.String(), nil)
	if err != nil {
		log.Err(err).Msg("error getting ending eth balance")
		return nil, err
	}
	tf.UserTrade.PostTradeEthBalance = ethBal
	tf.UserTrade.AddTxHash(accounts.Hash(tf.Tx.Hash()))
	return scInfo, err
}

func (u *UniswapClient) ExecFrontRunTradeStep(tf *TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	if u.DebugPrint {
		fmt.Println("executing front run trade")
	}
	return u.ExecSwap(tf.InitialPair, &tf.FrontRunTrade)
}

func (u *UniswapClient) ExecSandwichTradeStepTokenTransfer(tf *TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	_, _ = u.SandwichTradeGetAmountsOut(tf)
	ethBal, err := u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting pre trade eth balance")
		return nil, err
	}
	tf.SandwichTrade.PreTradeEthBalance = ethBal
	err = u.RouterApproveAndSend(ctx, &tf.SandwichTrade, tf.InitialPair.PairContractAddr)
	if err != nil {
		return nil, err
	}
	ethBal, err = u.Web3Client.GetBalance(ctx, u.Web3Client.PublicKey(), nil)
	if err != nil {
		log.Err(err).Msg("error getting post trade eth balance")
		return nil, err
	}
	tf.SandwichTrade.PostTradeEthBalance = ethBal
	return u.ExecSandwichTradeStep(tf)
}

func (u *UniswapClient) ExecSandwichTradeStep(tf *TradeExecutionFlow) (*web3_actions.SendContractTxPayload, error) {
	return u.ExecSwap(tf.InitialPair, &tf.SandwichTrade)
}

func (u *UniswapClient) FrontRunTradeGetAmountsOut(tf *TradeExecutionFlow) ([]*big.Int, error) {
	pathSlice := []string{tf.FrontRunTrade.AmountInAddr.String(), tf.FrontRunTrade.AmountOutAddr.String()}
	amountsOut, err := u.GetAmountsOut(tf.FrontRunTrade.AmountIn, pathSlice)
	if err != nil {
		return nil, err
	}
	amountsOutFirstPair := ConvertAmountsToBigIntSlice(amountsOut)
	if len(amountsOutFirstPair) != 2 {
		return nil, errors.New("amounts out not equal to expected")
	}
	if u.DebugPrint {
		fmt.Println("front run trade trade path", pathSlice[0], pathSlice[1])
		fmt.Println("front run trade expected amount in", tf.FrontRunTrade.AmountIn.String(), "amount out", tf.FrontRunTrade.AmountOut.String())
		fmt.Println("front run trade simulated amount in", amountsOutFirstPair[0].String(), "amount out", amountsOutFirstPair[1].String())
	}
	if tf.FrontRunTrade.AmountIn.String() != amountsOutFirstPair[0].String() {
		log.Warn().Msgf(fmt.Sprintf("amount in not equal to expected amount in %s, actual amount in: %s", tf.FrontRunTrade.AmountIn.String(), amountsOutFirstPair[0].String()))
		return amountsOutFirstPair, errors.New("amount in not equal to expected")
	}
	if tf.FrontRunTrade.AmountOut.String() != amountsOutFirstPair[1].String() {
		log.Warn().Msgf(fmt.Sprintf("amount out not equal to expected amount out %s, actual amount out: %s", tf.FrontRunTrade.AmountOut.String(), amountsOutFirstPair[1].String()))
		if u.DebugPrint {
			diff := new(big.Int).Sub(amountsOutFirstPair[1], tf.FrontRunTrade.AmountOut)
			fmt.Println("front run trade actual - expected ", diff.String())
		}
		return amountsOutFirstPair, errors.New("amount out not equal to expected")
	}
	tf.FrontRunTrade.SimulatedAmountOut = amountsOutFirstPair[1]
	return amountsOutFirstPair, err
}

func (u *UniswapClient) UserTradeGetAmountsOut(tf *TradeExecutionFlow) ([]*big.Int, error) {
	pathSlice := []string{tf.UserTrade.AmountInAddr.String(), tf.UserTrade.AmountOutAddr.String()}
	amountsOut, err := u.GetAmountsOut(tf.UserTrade.AmountIn, pathSlice)
	if err != nil {
		return nil, err
	}
	amountsOutFirstPair := ConvertAmountsToBigIntSlice(amountsOut)
	if len(amountsOutFirstPair) != 2 {
		return nil, errors.New("amounts out not equal to expected")
	}
	if u.DebugPrint {
		fmt.Println("user trade trade path", pathSlice[0], pathSlice[1])
		fmt.Println("user trade expected amount in", tf.UserTrade.AmountIn.String(), "amount out", tf.UserTrade.AmountOut.String())
		fmt.Println("user trade simulated amount in", amountsOutFirstPair[0].String(), "amount out", amountsOutFirstPair[1].String())
	}
	if tf.UserTrade.AmountIn.String() != amountsOutFirstPair[0].String() {
		log.Warn().Msgf(fmt.Sprintf("amount in not equal to expected amount in %s, actual amount in: %s", tf.UserTrade.AmountIn.String(), amountsOutFirstPair[0].String()))
		return amountsOutFirstPair, errors.New("amount in not equal to expected")
	}
	if tf.UserTrade.AmountOut.String() != amountsOutFirstPair[1].String() {
		log.Warn().Msgf(fmt.Sprintf("amount out not equal to expected amount out %s, actual amount out: %s", tf.UserTrade.AmountOut.String(), amountsOutFirstPair[1].String()))
		if u.DebugPrint {
			diff := new(big.Int).Sub(amountsOutFirstPair[1], tf.UserTrade.AmountOut)
			fmt.Println("user trade actual - expected ", diff.String())
		}
		return amountsOutFirstPair, errors.New("amount out not equal to expected")
	}
	tf.UserTrade.SimulatedAmountOut = amountsOutFirstPair[1]
	return amountsOutFirstPair, err
}

func (u *UniswapClient) SandwichTradeGetAmountsOut(tf *TradeExecutionFlow) ([]*big.Int, error) {
	pathSlice := []string{tf.SandwichTrade.AmountInAddr.String(), tf.SandwichTrade.AmountOutAddr.String()}
	amountsOut, err := u.GetAmountsOut(tf.SandwichTrade.AmountIn, pathSlice)
	if err != nil {
		return nil, err
	}
	amountsOutFirstPair := ConvertAmountsToBigIntSlice(amountsOut)
	if len(amountsOutFirstPair) != 2 {
		return nil, errors.New("amounts out not equal to expected")
	}
	if u.DebugPrint {
		fmt.Println("sandwich trade trade path", pathSlice[0], pathSlice[1])
		fmt.Println("sandwich trade expected amount in", tf.SandwichTrade.AmountIn.String(), "amount out", tf.SandwichTrade.AmountOut.String())
		fmt.Println("sandwich trade simulated amount in", amountsOutFirstPair[0].String(), "amount out", amountsOutFirstPair[1].String())
	}
	if tf.SandwichTrade.AmountIn.String() != amountsOutFirstPair[0].String() {
		log.Warn().Msgf(fmt.Sprintf("amount in not equal to expected amount in %s, actual amount in: %s", tf.UserTrade.AmountIn.String(), amountsOutFirstPair[0].String()))
		return amountsOutFirstPair, errors.New("amount in not equal to expected")
	}
	if tf.SandwichTrade.AmountOut.String() != amountsOutFirstPair[1].String() {
		log.Warn().Msgf(fmt.Sprintf("amount out not equal to expected amount out %s, actual amount out: %s", tf.UserTrade.AmountOut.String(), amountsOutFirstPair[1].String()))
		return amountsOutFirstPair, errors.New("amount out not equal to expected")
	}
	tf.SandwichTrade.SimulatedAmountOut = amountsOutFirstPair[1]
	return amountsOutFirstPair, err
}
