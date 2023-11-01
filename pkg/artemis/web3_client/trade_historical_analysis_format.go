package web3_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
)

func (u *UniswapClient) RunHistoricalTradeAnalysis(ctx context.Context, tfStr string) error {
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
	if u.Web3Client.IsAnvilNode == true {
		if u.Web3Client.Headers != nil {
			sid := u.Web3Client.Headers["Session-Lock-ID"]
			if sid != "" {
				defer u.EndHardHatSessionAndReset()
			}
		}
	}
	err = FilterNonActionTradeExecutionFlows(tfJSON)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	tf, err := tfJSON.ConvertToBigIntType()
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
	}
	_, err = u.CheckBlockRxAndNetworkReset(ctx, &tf)
	if err != nil {
		return u.MarkEndOfSimDueToErr(err)
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

func (u *UniswapClient) EndHardHatSessionAndReset() error {
	nodeInfo, err := u.Web3Client.GetNodeMetadata(ctx)
	if err != nil {
		return err
	}
	err = u.Web3Client.EndHardHatSessionReset(ctx, nodeInfo.ForkConfig.ForkUrl, 0)
	if err != nil {
		log.Err(err).Msg("UniswapClient: EndHardHatSessionReset")
		return err
	}
	return nil
}

// 17688038

func (u *UniswapClient) CheckBlockRxAndNetworkReset(ctx context.Context, tf *TradeExecutionFlow) (int, error) {
	artemis_test_cache.LiveTestNetwork.Dial()
	rx, err := artemis_test_cache.LiveTestNetwork.C.TransactionReceipt(ctx, tf.Tx.Hash())
	if err != nil {
		return -1, err
	}
	artemis_test_cache.LiveTestNetwork.Close()
	currentBlockStr := tf.CurrentBlockNumber.String()
	currentBlockNum, err := strconv.Atoi(currentBlockStr)
	if err != nil {
		return -1, err
	}
	u.TradeAnalysisReport.ArtemisBlockNumber = currentBlockNum
	u.TradeAnalysisReport.RxBlockNumber = int(rx.BlockNumber.Int64())
	if currentBlockNum >= int(rx.BlockNumber.Int64()) {
		err = fmt.Errorf("artmeis block number %d is greater than or equal to rx block number %d", currentBlockNum, int(rx.BlockNumber.Int64()))
		log.Err(err).Msg("CheckBlockRxAndNetworkReset: artemis block number is greater than or equal to rx block number")
		return -1, err
	}
	u.Web3Client.Dial()
	u.Web3Client.AddDefaultEthereumMainnetTableHeader()
	origInfo, err := u.Web3Client.GetNodeMetadata(ctx)
	if err != nil {
		log.Err(err).Msg("CheckBlockRxAndNetworkReset: error getting node metadata")
		return -1, err
	}
	u.Web3Client.Close()
	u.Web3Client.Dial()
	defer u.Web3Client.Close()
	err = u.Web3Client.ResetNetwork(ctx, origInfo.ForkConfig.ForkUrl, currentBlockNum)
	if err != nil {
		log.Err(err).Msg("CheckBlockRxAndNetworkReset: error resetting network")
		return -1, err
	}
	u.Web3Client.Close()
	for i := 0; i < 10; i++ {
		u.Web3Client.Dial()
		nodeInfo, rerr := u.Web3Client.GetNodeMetadata(ctx)
		if rerr != nil {
			return -1, err
		}
		u.Web3Client.Close()
		if nodeInfo.ForkConfig.ForkUrl != origInfo.ForkConfig.ForkUrl {
			return -1, fmt.Errorf("CheckBlockRxAndNetworkReset: live network fork url %s is not equal to initial fork url %s", nodeInfo.ForkConfig.ForkUrl, nodeInfo.ForkConfig.ForkUrl)
		}
		if nodeInfo.ForkConfig.ForkBlockNumber != currentBlockNum {
			fmt.Println("initForkUrl1", origInfo.ForkConfig.ForkUrl, "CurrentBlockNumber", origInfo.CurrentBlockNumber.ToInt().String(), "ForkBlockNumber", origInfo.ForkConfig.ForkBlockNumber)
			fmt.Println("initForkUrl2", nodeInfo.ForkConfig.ForkUrl, "CurrentBlockNumber", nodeInfo.CurrentBlockNumber.ToInt().String(), "ForkBlockNumber", nodeInfo.ForkConfig.ForkBlockNumber)
		} else {
			return currentBlockNum, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return -1, errors.New("CheckBlockRxAndNetworkReset: could not reset network")
}

func (u *UniswapClient) CheckExpectedReserves(tf *TradeExecutionFlow) error {
	if tf.InitialPair == nil {
		return nil
	}
	// todo, do v3 pairs
	simPair := tf.InitialPair
	err := uniswap_pricing.GetPairContractPrices(ctx, tf.CurrentBlockNumber.Uint64(), u.Web3Client.Web3Actions, simPair)
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
