package artemis_trade_debugger

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
)

/*
type TradeExecutionFlow struct {
	CurrentBlockNumber *big.Int                           `json:"currentBlockNumber"`
	Tx                 *types.Transaction                 `json:"tx"`
	Trade              Trade                              `json:"trade"`
	InitialPair        *uniswap_pricing.UniswapV2Pair     `json:"initialPair,omitempty"`
	InitialPairV3      *uniswap_pricing.UniswapV3Pair     `json:"initialPairV3,omitempty"`
	FrontRunTrade      artemis_trading_types.TradeOutcome `json:"frontRunTrade"`
	UserTrade          artemis_trading_types.TradeOutcome `json:"userTrade"`
	SandwichTrade      artemis_trading_types.TradeOutcome `json:"sandwichTrade"`
	SandwichPrediction SandwichTradePrediction            `json:"sandwichPrediction"`
}
*/

// 0x80ae3cc1748c10f42e591783001817b8a56b188eb1867282e396a8d99d583d00

// {"level":"warn","txHash":"0x43dd0f388b41b536e50bc25de1238aa46b3e341bc3d98b26c94fbed184537590",
// "tradeMethod":"V2_SWAP_EXACT_IN","toAddr":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","time":1689974271,
// "message":"dat: ApplyMaxTransferTax, tokenOne and tokenTwo are zero address"}

func (t *ArtemisTradeDebuggerTestSuite) TestReplayer() {
	// 0x925dd1373fea0f4537e9670dc984a5c0640da81142269e8eff6840d8caaea6f4
	txHash := "0xf1ed952cff38e1941ba947a0bf5ee12e6d70bfbdbc8f3b8ebbad99372dd1ac4f"
	t.td.dat.GetSimUniswapClient().Web3Client.AddSessionLockHeader(txHash)
	err := t.td.Replay(ctx, txHash, true)
	t.NoError(err)
}

func (t *ArtemisTradeDebuggerTestSuite) TestReplayerWithoutTax() {
	// 0x925dd1373fea0f4537e9670dc984a5c0640da81142269e8eff6840d8caaea6f4
	txHash := "0xf1ed952cff38e1941ba947a0bf5ee12e6d70bfbdbc8f3b8ebbad99372dd1ac4f"
	t.td.dat.GetSimUniswapClient().Web3Client.AddSessionLockHeader(txHash)
	err := t.td.ReplayWithoutTax(ctx, txHash, true)
	t.NoError(err)
}

func (t *ArtemisTradeDebuggerTestSuite) TestReplayerBulk() {
	txs, err := artemis_mev_models.SelectReplayEthMevMempoolTxByTxHash(ctx)
	t.NoError(err)
	for _, txMem := range txs {
		txHash := txMem.EthMevTxAnalysis.TxHash
		err = t.td.Replay(ctx, txHash, true)
		t.NoError(err)
	}
}

// artemis_mev_models
/*
0x4a9c05ef46a2a0f4d36577bd38e37502245448a1b52da9c73ca59af37059f89e
profitToken 0x0359181dCE76bAD4d3f851b3356FdD7b82A41B14
expectedProfit 7807642577146113
actualProfit 8695022423393079

0x925dd1373fea0f4537e9670dc984a5c0640da81142269e8eff6840d8caaea6f4
profitToken 0x511686014F39F487E5CDd5C37B4b37606B795ae3
expectedProfit 6635478652156427361470498
actualProfit 6788708842908401012256112
*/

// 0x58282b7b489ae24a75e7b49b68f1360d95374e00a4dbc58c3aaea3329c4e8aca
func (t *ArtemisTradeDebuggerTestSuite) TestReadRx() {
	artemis_test_cache.LiveTestNetwork.Dial()
	defer artemis_test_cache.LiveTestNetwork.Close()
	txHash := "0x58282b7b489ae24a75e7b49b68f1360d95374e00a4dbc58c3aaea3329c4e8aca"
	rx, err := artemis_test_cache.LiveTestNetwork.C.TransactionReceipt(ctx, common.HexToHash(txHash))
	t.NoError(err)

	fmt.Println(rx.ContractAddress.String())
	fmt.Println(rx.BlockNumber.String())
	fmt.Println(rx.Status)
	fmt.Println(rx.GasUsed)
}
