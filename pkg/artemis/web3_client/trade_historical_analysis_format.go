package web3_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
)

func UnmarshalTradeExecutionFlow(tfStr string) (TradeExecutionFlow, error) {
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
	TxHash             string `json:"tx_hash"`
	TradeMethod        string `json:"trade_method"`
	ArtemisBlockNumber int    `json:"artemis_block_number"`
	RxBlockNumber      int    `json:"rx_block_number"`

	GasReport          `json:"gas_report"`
	TradeFailureReport `json:"trade_failure_report"`
	SimulationResults  `json:"simulation_results"`
}

func (t *TradeAnalysisReport) SaveResultsInDb(ctx context.Context) error {
	b, err := json.Marshal(t)
	if err != nil {
		log.Err(err).Msg("error marshalling trade analysis report")
		return err
	}
	txAnalysis := artemis_autogen_bases.EthMevTxAnalysis{
		GasUsedWei:              t.TotalGasUsed,
		Metadata:                string(b),
		TxHash:                  t.TxHash,
		TradeMethod:             t.TradeMethod,
		EndReason:               t.EndReason,
		AmountIn:                t.AmountIn,
		AmountOutAddr:           t.AmountOutAddr,
		ExpectedProfitAmountOut: t.ExpectedProfitAmountOut,
		RxBlockNumber:           t.RxBlockNumber,
		AmountInAddr:            t.AmountInAddr,
		ActualProfitAmountOut:   t.AmountOut,
	}
	if strings.HasSuffix(t.EndReason, "quiknode.com") {
		return errors.New("rate limit error")
	}
	err = artemis_validator_service_groups_models.InsertEthMevTxAnalysis(ctx, txAnalysis)
	if err != nil {
		log.Err(err).Msg("error inserting into eth_mev_tx_analysis")
		return err
	}
	return nil
}
func (t *TradeAnalysisReport) PrintResults() {
	if t.EndReason == "trade failed due to invalid amount in" {
		return
	}
	fmt.Println("Trade Method:", t.TradeMethod)
	fmt.Println("Artemis Block Number:", t.ArtemisBlockNumber)
	fmt.Println("Rx Block Number:", t.RxBlockNumber)

	if t.EndStage == "success" {
		fmt.Println("Starting Token Addr:", t.AmountInAddr)
		fmt.Println("Profit Token Addr:", t.AmountOutAddr)
		fmt.Println("Actual Profit:", t.AmountOut)
		fmt.Println("Expected Profit:", t.ExpectedProfitAmountOut)
		fmt.Println("Total Gas Used:", t.TotalGasUsed)
	} else {
		fmt.Println("End Reason:", t.EndReason)
		fmt.Println("End Stage:", t.EndStage)
	}
}

type SimulationResults struct {
	AmountInAddr            string `json:"amount_in_addr"`
	AmountIn                string `json:"amount_in"`
	AmountOutAddr           string `json:"amount_out_addr"`
	AmountOut               string `json:"amount_out"`
	ExpectedProfitAmountOut string `json:"expected_profit_amount_out"`
}

type GasReport struct {
	TotalGasUsed         string `json:"total_gas_used"`
	FrontRunGasUsed      string `json:"front_run_gas_used"`
	SandwichTradeGasUsed string `json:"sandwich_trade_gas_used"`
}

type TradeFailureReport struct {
	EndReason string `json:"end_reason"`
	EndStage  string `json:"end_stage"`
}

func (u *UniswapClient) RunHistoricalTradeAnalysis(ctx context.Context, tfStr string, liveNetworkClient Web3Client) error {
	u.TradeAnalysisReport = &TradeAnalysisReport{}
	tfJSON, err := UnmarshalTradeExecutionFlow(tfStr)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	u.TradeAnalysisReport.TxHash = tfJSON.Tx.Hash().String()
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
	u.TradeAnalysisReport.AmountIn = tf.FrontRunTrade.AmountIn.String()
	u.TradeAnalysisReport.AmountInAddr = tf.FrontRunTrade.AmountInAddr.String()
	u.TradeAnalysisReport.AmountOutAddr = tf.SandwichTrade.AmountOutAddr.String()
	err = u.SimFullSandwichTrade(&tf)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	return nil
}

func (u *UniswapClient) CheckBlockRxAndNetworkReset(tf *TradeExecutionFlowInBigInt, liveNetworkClient Web3Client) (int, error) {
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

func (u *UniswapClient) CheckExpectedReserves(tf *TradeExecutionFlowInBigInt) error {
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

func (t *TradeAnalysisReport) MarkEndOfSimDueToErr(err error) error {
	// mark end of test
	if err != nil {
		t.EndReason = err.Error()
	} else {
		t.EndReason = "success"
	}
	return t.SaveResultsInDb(ctx)
}
