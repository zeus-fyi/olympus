package artemis_trade_debugger

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/dynamic_secrets"
	artemis_trading_auxiliary "github.com/zeus-fyi/olympus/pkg/artemis/trading/auxiliary"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_multicall "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/multicall"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
	"github.com/zeus-fyi/olympus/pkg/athena"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
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

const (
	AccountAddr = "0x000000641e80A183c8B736141cbE313E136bc8c6"
)

func (t *ArtemisTradeDebuggerTestSuite) TestEthCall() {
	// CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	// func (ec *Client) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	//PendingCallContract
	//CallContract
	s3client, err := s3base.NewConnS3ClientWithStaticCreds(ctx, t.Tc.LocalS3SpacesKey, t.Tc.LocalS3SpacesSecret)
	t.Require().Nil(err)
	athena.AthenaS3Manager = s3client

	age := encryption.NewAge(t.Tc.LocalAgePkey, t.Tc.LocalAgePubkey)
	acc := initTradingAccount(ctx, age)

	tx, err := GetMevMempoolTxTradeFlow(ctx, "0xf0824882281da217321b75c40a4f24e9d3c88443d2f3741c2e43b971cad624f6")
	t.Assert().Nil(err)
	t.Assert().NotNil(tx)

	txInt, err := tx.ConvertToBigIntType()
	t.Assert().Nil(err)
	t.Assert().NotNil(txInt)

	wc := web3_client.NewWeb3Client("https://iris.zeus.fyi/v1/router", &acc)
	//wc = web3_client.NewWeb3ClientFakeSigner("https://iris.zeus.fyi/v1/router")

	wc.AddDefaultEthereumMainnetTableHeader()
	wc.AddBearerToken(t.Tc.ProductionLocalTemporalBearerToken)
	wc.Dial()
	defer wc.Close()
	fmt.Println(txInt.FrontRunTrade.AmountIn.String())

	fmt.Println(txInt.FrontRunTrade.AmountInAddr.String())
	fmt.Println(txInt.FrontRunTrade.AmountOutAddr.String())

	txInt.FrontRunTrade.AmountOut = new(big.Int).SetInt64(0)
	ur, _, err := artemis_trading_auxiliary.GenerateTradeV2SwapFromTokenToToken(ctx, wc, nil, &txInt.FrontRunTrade)
	t.Assert().Nil(err)
	urAbi := artemis_oly_contract_abis.MustLoadNewUniversalRouterAbi()
	encParams, err := ur.EncodeCommands(ctx, urAbi)
	t.Assert().Nil(err)
	t.Assert().NotNil(encParams)

	data, err := artemis_trading_auxiliary.GetUniswapUniversalRouterAbiPayload(ctx, encParams)
	t.Assert().Nil(err)
	t.Assert().NotNil(data)

	to := common.HexToAddress(artemis_trading_constants.UniswapUniversalRouterAddressNew)
	msg := ethereum.CallMsg{
		From:      common.HexToAddress(AccountAddr),
		To:        &to,
		Gas:       data.GasLimit,
		GasPrice:  data.GasPrice,
		GasFeeCap: data.GasFeeCap,
		GasTipCap: data.GasTipCap,
		Data:      data.Data,
	}
	resp, err := wc.C.CallContract(ctx, msg, nil)
	t.Assert().Nil(err)
	t.Assert().NotNil(resp)
	if err != nil {
		fmt.Println(err.Error())
	}

	m3 := []artemis_multicall.MultiCallElement{
		{
			Name: "balanceOf",
			Call: artemis_multicall.Call{
				Target:       common.HexToAddress(artemis_trading_constants.WETH9ContractAddress),
				AllowFailure: false,
				Data:         nil,
			},
			AbiFile:       artemis_oly_contract_abis.MustLoadERC20Abi(),
			DecodedInputs: []interface{}{common.HexToAddress(AccountAddr)},
		},
		{
			Name: data.MethodName,
			Call: artemis_multicall.Call{
				Target:       common.HexToAddress(to.String()),
				AllowFailure: true,
				Data:         data.Data,
			},
			AbiFile:       urAbi,
			DecodedInputs: data.Params,
		},
		{
			Name: "balanceOf",
			Call: artemis_multicall.Call{
				Target:       common.HexToAddress(artemis_trading_constants.WETH9ContractAddress),
				AllowFailure: false,
				Data:         nil,
			},
			AbiFile:       artemis_oly_contract_abis.MustLoadERC20Abi(),
			DecodedInputs: []interface{}{common.HexToAddress(AccountAddr)},
		},
	}
	/*
			{
			Name: data.MethodName,
			Call: artemis_multicall.Call{
				Target:       common.HexToAddress(artemis_trading_constants.UniswapUniversalRouterAddressNew),
				AllowFailure: true,
				Data:         data.Data,
			},
			AbiFile:       urAbi,
			DecodedInputs: data.Params,
		},
	*/
	payload, err := artemis_multicall.UrCreateMulticall3Payload(ctx, m3, data)
	t.Require().Nil(err)
	t.Assert().NotNil(payload)

	err = payload.GenerateBinDataFromParamsAbi(ctx)
	t.Require().Nil(err)

	re, err := wc.CallConstantFunction(ctx, &payload)
	t.Assert().Nil(err)
	t.Assert().NotNil(re)
	fmt.Println(re)

	out, err := artemis_multicall.UnpackMultiCall(ctx, re, m3)
	t.Assert().Nil(err)
	t.Assert().NotNil(out)
	fmt.Println(out)
}

// 0x08c379a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000174d756c746963616c6c333a2063616c6c206661696c6564000000000000000000
// TODO, need to multicall this i fucking guess. fuck you ethereum

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

// InitTradingAccount pubkey 0x000025e60C7ff32a3470be7FE3ed1666b0E326e2
func initTradingAccount(ctx context.Context, age encryption.Age) accounts.Account {
	p := filepaths.Path{
		DirIn:  "keygen",
		DirOut: "keygen",
		FnIn:   "key-6.txt.age",
	}
	r, err := dynamic_secrets.ReadAddress(ctx, p, athena.AthenaS3Manager, age)
	if err != nil {
		panic(err)
	}
	acc, err := dynamic_secrets.GetAccount(r)
	if err != nil {
		panic(err)
	}
	return acc
}
