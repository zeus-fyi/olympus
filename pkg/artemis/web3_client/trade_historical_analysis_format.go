package web3_client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

func UnmarshallTradeExecutionFlow(tfStr string) (TradeExecutionFlow, error) {
	tf := TradeExecutionFlow{}
	by := []byte(tfStr)
	berr := json.Unmarshal(by, &tf)
	if berr != nil {
		return tf, berr
	}
	return tf, nil
}

func FilterNonActionTradeExecutionFlows(tf TradeExecutionFlow) error {
	if tf.FrontRunTrade.AmountIn == "0" || tf.FrontRunTrade.AmountIn == "" {
		return fmt.Errorf("trade failed due to invalid amount in")
	}
	return nil
}

type TradeAnalysisReport struct {
	TradeMethod        string
	ArtemisBlockNumber int
	RxBlockNumber      int

	GasReport
	TradeFailureReport
	SimulationResults
}

type SimulationResults struct {
	StartingTokenAddr string
	ProfitTokenAddr   string
	ActualProfit      string
	ExpectedProfit    string
}

type GasReport struct {
	TotalGasUsed         string
	FrontRunGasUsed      string
	SandwichTradeGasUsed string
}

type TradeFailureReport struct {
	EndReason string
	EndStage  string
}

func (u *UniswapV2Client) RunHistoricalTradeAnalysis(ctx context.Context, tfStr string, liveNetworkClient Web3Client) error {
	u.TradeAnalysisReport = &TradeAnalysisReport{}
	tfJSON, err := UnmarshallTradeExecutionFlow(tfStr)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	u.TradeAnalysisReport.TradeMethod = tfJSON.Trade.TradeMethod
	err = FilterNonActionTradeExecutionFlows(tfJSON)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	tf := tfJSON.ConvertToBigIntType()
	artemisBlockNum, err := u.CheckBlockRxAndNetworkReset(&tf, liveNetworkClient)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	err = u.Web3Client.HardHatResetNetwork(ctx, liveNetworkClient.NodeURL, artemisBlockNum)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	err = u.CheckExpectedReserves(&tf)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	u.TradeAnalysisReport.StartingTokenAddr = tf.FrontRunTrade.AmountInAddr.String()
	err = u.SimFullSandwichTrade(&tf)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	u.EndReason = "success"
	return nil
}

func (u *UniswapV2Client) CheckBlockRxAndNetworkReset(tf *TradeExecutionFlowInBigInt, liveNetworkClient Web3Client) (int, error) {
	rx, err := liveNetworkClient.GetTxReceipt(ctx, tf.Tx.Hash())
	if err != nil {
		return -1, err
	}
	currentBlockStr := tf.CurrentBlockNumber.String()
	currentBlockNum, err := strconv.Atoi(currentBlockStr)
	if err != nil {
		return -1, err
	}
	u.TradeAnalysisReport.ArtemisBlockNumber = currentBlockNum
	u.TradeAnalysisReport.RxBlockNumber = int(rx.BlockNumber.Int64())
	if currentBlockNum >= int(rx.BlockNumber.Int64()) {
		return -1, fmt.Errorf("artmeis block number %d is greater than or equal to rx block number %d", currentBlockNum, int(rx.BlockNumber.Int64()))
	}
	return currentBlockNum, nil
}

func (u *UniswapV2Client) CheckExpectedReserves(tf *TradeExecutionFlowInBigInt) error {
	pairAddr := tf.InitialPair.PairContractAddr
	simPair, err := u.GetPairContractPrices(ctx, pairAddr)
	if err != nil {
		return err
	}
	if tf.InitialPair.Reserve1.String() != simPair.Reserve1.String() && tf.InitialPair.Reserve0.String() != simPair.Reserve0.String() {
		return fmt.Errorf("reserve mismatch")
	}
	if tf.InitialPair.Reserve0.String() != simPair.Reserve0.String() {
		return fmt.Errorf("reserve0 mismatch")
	}
	if tf.InitialPair.Reserve1.String() != simPair.Reserve1.String() {
		return fmt.Errorf("reserve1 mismatch")
	}
	return nil
}

func (u *TradeAnalysisReport) MarkEndOfSimDueToErr(err error) error {
	// mark end of test
	u.EndReason = err.Error()
	return err
}
