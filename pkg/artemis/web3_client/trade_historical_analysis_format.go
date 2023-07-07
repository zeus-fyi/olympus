package web3_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
)

func (u *UniswapClient) RunHistoricalTradeAnalysis(ctx context.Context, tfStr string, liveNetworkClient Web3Client) error {
	u.TradeAnalysisReport = &TradeAnalysisReport{}
	tfJSON, err := UnmarshalTradeExecutionFlow(tfStr)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	pairAddr := ""
	if tfJSON.InitialPairV3 != nil {
		pairAddr = tfJSON.InitialPairV3.PoolAddress
	} else if tfJSON.InitialPair != nil {
		pairAddr = tfJSON.InitialPair.PairContractAddr
	}
	u.TradeAnalysisReport.PairAddress = pairAddr
	u.TradeAnalysisReport.TxHash = tfJSON.Tx.Hash
	u.TradeAnalysisReport.TradeMethod = tfJSON.Trade.TradeMethod
	u.TradeAnalysisReport.AmountIn = tfJSON.FrontRunTrade.AmountIn
	u.TradeAnalysisReport.AmountInAddr = tfJSON.UserTrade.AmountInAddr.String()
	u.TradeAnalysisReport.AmountOutAddr = tfJSON.UserTrade.AmountOutAddr.String()
	u.Web3Client.AddSessionLockHeader(tfJSON.Tx.Hash)

	err = FilterNonActionTradeExecutionFlows(tfJSON)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	tf := tfJSON.ConvertToBigIntType()
	artemisBlockNum, err := u.CheckBlockRxAndNetworkReset(ctx, &tf, &liveNetworkClient)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	if u.Web3Client.IsAnvilNode == true {
		if u.Web3Client.Headers != nil {
			sid := u.Web3Client.Headers["Session-Lock-ID"]
			if sid != "" {
				defer u.Web3Client.EndHardHatSessionReset(ctx, liveNetworkClient.NodeURL, artemisBlockNum)
			}
		}
	}
	err = u.CheckExpectedReserves(&tf)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	err = u.SimFullSandwichTrade(&tf)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	return nil
}

func FilterNonActionTradeExecutionFlows(tf TradeExecutionFlowJSON) error {
	if tf.FrontRunTrade.AmountIn == "0" || tf.FrontRunTrade.AmountIn == "" {
		return fmt.Errorf("trade failed due to invalid amount in")
	}
	return nil
}

type TradeAnalysisReport struct {
	TxHash             string `json:"txHash"`
	TradeMethod        string `json:"tradeMethod"`
	ArtemisBlockNumber int    `json:"artemisBlockNumber"`
	RxBlockNumber      int    `json:"rxBlockNumber"`
	PairAddress        string `json:"pairAddress,omitempty"`
	GasReport          `json:"gasReport"`
	TradeFailureReport `json:"tradeFailureReport"`
	SimulationResults  `json:"simulationResults"`
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
		PairAddress:             t.PairAddress,
	}
	if strings.HasSuffix(t.EndReason, "quiknode.com") {
		return errors.New("rate limit error")
	}
	err = artemis_mev_models.InsertEthMevTxAnalysis(ctx, txAnalysis)
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
	fmt.Println("Tx Hash:", t.TxHash)
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

func (u *UniswapClient) CheckBlockRxAndNetworkReset(ctx context.Context, tf *TradeExecutionFlow, liveNetworkClient *Web3Client) (int, error) {
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
	u.Web3Client.Dial()
	u.Web3Client.Close()
	nodeInfo, err := u.Web3Client.GetNodeMetadata(ctx)
	if err != nil {
		return -1, err
	}
	liveNetwork := nodeInfo.ForkConfig.ForkUrl
	liveNetworkClient.NodeURL = liveNetwork
	err = u.Web3Client.ResetNetwork(ctx, liveNetworkClient.NodeURL, currentBlockNum)
	if err != nil {
		return -1, err
	}

	nodeInfo, err = u.Web3Client.GetNodeMetadata(ctx)
	if err != nil {
		return -1, err
	}
	if nodeInfo.ForkConfig.ForkBlockNumber != currentBlockNum {
		return -1, fmt.Errorf("CheckBlockRxAndNetworkReset: live network fork block number %d is not equal to current block number %d", nodeInfo.ForkConfig.ForkBlockNumber, currentBlockNum)
	}
	return currentBlockNum, nil
}

func (u *UniswapClient) CheckExpectedReserves(tf *TradeExecutionFlow) error {
	if tf.InitialPair == nil {
		return nil
	}
	// todo, do v3 pairs
	simPair := tf.InitialPair
	err := uniswap_pricing.GetPairContractPrices(ctx, u.Web3Client.Web3Actions, simPair)
	if err != nil {
		log.Err(err).Msg("error getting pair contract prices")
		return err
	}
	if tf.InitialPair.Reserve1.String() != simPair.Reserve1.String() && tf.InitialPair.Reserve0.String() != simPair.Reserve0.String() {
		fmt.Println("tf.InitialPair.Reserve0", tf.InitialPair.Reserve0.String(), simPair.Reserve0.String(), "simPair.Reserve0")
		fmt.Println("tf.InitialPair.Reserve1", tf.InitialPair.Reserve1.String(), simPair.Reserve1.String(), "simPair.Reserve1")
		return fmt.Errorf("reserve mismatch")
	}
	if tf.InitialPair.Reserve0.String() != simPair.Reserve0.String() {
		fmt.Println("tf.InitialPair.Reserve0", tf.InitialPair.Reserve0.String(), simPair.Reserve0.String(), "simPair.Reserve0")
		return fmt.Errorf("reserve0 mismatch")
	}
	if tf.InitialPair.Reserve1.String() != simPair.Reserve1.String() {
		fmt.Println("tf.InitialPair.Reserve1", tf.InitialPair.Reserve1.String(), simPair.Reserve1.String(), "simPair.Reserve1")
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
