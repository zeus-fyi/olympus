package web3_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/gochain/gochain/v4/accounts/abi"
	"github.com/gochain/gochain/v4/common"
	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/signing_automation/ethereum"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

const (
	UniswapV2FactoryAddress = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
	UniswapV2RouterAddress  = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"

	addLiquidity                 = "addLiquidity"
	addLiquidityETH              = "addLiquidityETH"
	removeLiquidity              = "removeLiquidity"
	removeLiquidityETH           = "removeLiquidityETH"
	removeLiquidityWithPermit    = "removeLiquidityWithPermit"
	removeLiquidityETHWithPermit = "removeLiquidityETHWithPermit"
	swapExactTokensForTokens     = "swapExactTokensForTokens"
	swapTokensForExactTokens     = "swapTokensForExactTokens"
	swapExactETHForTokens        = "swapExactETHForTokens"
	swapTokensForExactETH        = "swapTokensForExactETH"
	swapExactTokensForETH        = "swapExactTokensForETH"
	swapETHForExactTokens        = "swapETHForExactTokens"
)

/*
https://docs.uniswap.org/contracts/v2/concepts/advanced-topics/fees
There is a 0.3% fee for swapping tokens. This fee is split by liquidity providers proportional to their contribution to liquidity reserves.
*/

// TODO https://docs.uniswap.org/contracts/v2/reference/smart-contracts/router-02

type UniswapV2Client struct {
	Web3Client               Web3Client
	FactorySmartContractAddr string
	PairAbi                  *abi.ABI
	ERC20Abi                 *abi.ABI
	FactoryAbi               *abi.ABI
	PrintOn                  bool
	PrintLocal               bool
	MevSmartContractTxMap
	Path        filepaths.Path
	BlockNumber *big.Int

	SwapExactTokensForTokensParamsSlice []SwapExactTokensForTokensParams
	SwapTokensForExactTokensParamsSlice []SwapTokensForExactTokensParams
	SwapExactETHForTokensParamsSlice    []SwapExactETHForTokensParams
	SwapTokensForExactETHParamsSlice    []SwapTokensForExactETHParams
	SwapExactTokensForETHParamsSlice    []SwapExactTokensForETHParams
	SwapETHForExactTokensParamsSlice    []SwapETHForExactTokensParams
}

func InitUniswapV2Client(ctx context.Context, w Web3Client) UniswapV2Client {
	abiFile, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.UniswapV2RouterABI))
	if err != nil {
		panic(err)
	}
	erc20AbiFile, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.ERC20ABI))
	if err != nil {
		panic(err)
	}
	factoryAbiFile, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.UniswapV2FactoryABI))
	if err != nil {
		panic(err)
	}
	pairAbiFile, err := signing_automation_ethereum.ReadAbi(ctx, strings.NewReader(artemis_oly_contract_abis.UniswapV2PairAbi))
	if err != nil {
		panic(err)
	}
	f := strings_filter.FilterOpts{
		DoesNotStartWithThese: nil,
		StartsWithThese:       []string{"swap"},
		Contains:              "",
		DoesNotInclude:        []string{"supportingFeeOnTransferTokens"},
	}
	return UniswapV2Client{
		Web3Client:               w,
		FactorySmartContractAddr: UniswapV2FactoryAddress,
		FactoryAbi:               factoryAbiFile,
		ERC20Abi:                 erc20AbiFile,
		PairAbi:                  pairAbiFile,
		MevSmartContractTxMap: MevSmartContractTxMap{
			SmartContractAddr: UniswapV2RouterAddress,
			Abi:               abiFile,
			MethodTxMap:       map[string]MevTx{},
			Txs:               []MevTx{},
			Filter:            &f,
		},
		SwapExactTokensForTokensParamsSlice: []SwapExactTokensForTokensParams{},
		SwapTokensForExactTokensParamsSlice: []SwapTokensForExactTokensParams{},
		SwapExactETHForTokensParamsSlice:    []SwapExactETHForTokensParams{},
		SwapTokensForExactETHParamsSlice:    []SwapTokensForExactETHParams{},
		SwapExactTokensForETHParamsSlice:    []SwapExactTokensForETHParams{},
		SwapETHForExactTokensParamsSlice:    []SwapETHForExactTokensParams{},
	}
}

func (u *UniswapV2Client) GetAllTradeMethods() []string {
	return []string{
		addLiquidity,
		addLiquidityETH,
		removeLiquidity,
		removeLiquidityETH,
		removeLiquidityWithPermit,
		removeLiquidityETHWithPermit,
		swapExactTokensForTokens,
		swapTokensForExactTokens,
		swapExactETHForTokens,
		swapTokensForExactETH,
		swapExactTokensForETH,
		swapETHForExactTokens,
	}
}

func (u *UniswapV2Client) ProcessTxs(ctx context.Context) {
	bn, err := u.Web3Client.GetHeadBlockHeight(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to get block height")
		return
	}
	u.BlockNumber = bn
	count := 0
	for methodName, tx := range u.MethodTxMap {
		switch methodName {
		case addLiquidity:
			//u.AddLiquidity(tx.Args)
		case addLiquidityETH:
			// payable
			//u.AddLiquidityETH(tx.Args)
			if tx.Tx.Value == nil {
				continue
			}
		case removeLiquidity:
			//u.RemoveLiquidity(tx.Args)
		case removeLiquidityETH:
			//u.RemoveLiquidityETH(tx.Args)
		case removeLiquidityWithPermit:
			//u.RemoveLiquidityWithPermit(tx.Args)
		case removeLiquidityETHWithPermit:
			//u.RemoveLiquidityETHWithPermit(tx.Args)
		case swapExactTokensForTokens:
			count++
			u.SwapExactTokensForTokens(tx, tx.Args)
		case swapTokensForExactTokens:
			count++
			u.SwapTokensForExactTokens(tx, tx.Args)
		case swapExactETHForTokens:
			// payable
			count++
			if tx.Tx.Value == nil {
				continue
			}
			u.SwapExactETHForTokens(tx, tx.Args, tx.Tx.Value.ToInt())
		case swapTokensForExactETH:
			count++
			u.SwapTokensForExactETH(tx, tx.Args)
		case swapExactTokensForETH:
			count++
			u.SwapExactTokensForETH(tx, tx.Args)
		case swapETHForExactTokens:
			// payable
			count++
			if tx.Tx.Value == nil {
				continue
			}
			u.SwapETHForExactTokens(tx, tx.Args, tx.Tx.Value.ToInt())
		}
	}
	fmt.Println("totalFilteredCount:", count)
}

func (u *UniswapV2Client) PrintTradeSummaries(tx MevTx, tf TradeExecutionFlow, pair UniswapV2Pair, tokenAddr string, amount, amountMin *big.Int) {
	tf.Tx = tx.Tx
	tf.CurrentBlockNumber = u.BlockNumber
	expectedOut, err := pair.GetQuoteUsingTokenAddr(tokenAddr, amount)
	if err != nil {
		fmt.Println("GetQuoteUsingTokenAddr", err)
		return
	}
	diff := new(big.Int).Sub(expectedOut, amountMin)
	purchasedTokenAddr := pair.GetOppositeToken(tokenAddr).String()
	fmt.Printf("Token0 Address: %s Token0 Reserve: %s,\nToken1 Address %s, Token1 Reserve: %s\n", pair.Token0.String(), pair.Reserve0.String(), pair.Token1.String(), pair.Reserve1.String())
	fmt.Printf("Expected amount %s %s token from trade at current rate \n", expectedOut.String(), purchasedTokenAddr)
	fmt.Printf("Amount minimum %s %s token needed from trade \n", amountMin.String(), purchasedTokenAddr)
	bn, err := u.Web3Client.GetHeadBlockHeight(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("failed to get block height")
		return
	}
	if u.BlockNumber.String() != bn.String() {
		log.Info().Interface("currentBlockNumber", bn.String()).Interface("startingBlockNumber", u.BlockNumber.String()).Msg("block number transition exiting due to stale data")
		return
	}
	if diff.Cmp(big.NewInt(0)) == 1 {
		fmt.Printf("Positive difference between expected and minimum amount is %s %s token \n", diff.String(), tokenAddr)
		b, berr := json.MarshalIndent(tf, "", "  ")
		if berr != nil {
			return
		}
		if u.PrintLocal {
			u.Path.FnOut = fmt.Sprintf("%s-%d.json", tf.TradeMethod, u.BlockNumber)
			err = u.Path.WriteToFileOutPath(b)
			if err != nil {
				return
			}
		}
		if tx.Tx.Nonce == nil {
			fmt.Printf("tx.Tx.Nonce is nil")
			return
		}
		btf, berr := json.Marshal(tf)
		if berr != nil {
			return
		}
		b, berr = json.Marshal(tf.Tx)
		if berr != nil {
			return
		}
		txMempool := artemis_autogen_bases.EthMempoolMevTx{
			ProtocolNetworkID: hestia_req_types.EthereumMainnetProtocolNetworkID,
			Tx:                string(b),
			TxFlowPrediction:  string(btf),
			TxHash:            tx.Tx.Hash.String(),
			Nonce:             int(*tx.Tx.Nonce),
			From:              tx.Tx.From.String(),
			To:                tx.Tx.To.String(),
			BlockNumber:       int(u.BlockNumber.Int64()),
		}
		err = artemis_validator_service_groups_models.InsertMempoolTx(ctx, txMempool)
		if err != nil {
			fmt.Printf("InsertMempoolTx err: %s", err)
			return
		}
	} else {
		fmt.Printf("Negative difference between expected and minimum amount is %s %s token \n", diff.String(), tokenAddr)
	}
	if amountMin.Cmp(big.NewInt(0)) == 0 {
		fmt.Printf("Amount minimum is 0, so no trade will be executed \n")
		return
	}
	slippage := new(big.Int).Mul(diff, big.NewInt(100))
	slippagePercent := new(big.Int).Div(slippage, amountMin)
	fmt.Printf("Slippage is %s %% \n", slippagePercent.String())
	fmt.Printf("Buy %s %s token for %s %s token \n\n", expectedOut.String(), pair.GetOppositeToken(tokenAddr).String(), amount.String(), tokenAddr)
	return
}

func (u *UniswapV2Client) SwapExactTokensForTokens(tx MevTx, args map[string]interface{}) {
	amountIn, err := ParseBigInt(args["amountIn"])
	if err != nil {
		return
	}
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		return
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		return
	}
	to, err := ConvertToAddress(args["to"])
	if err != nil {
		return
	}
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		return
	}
	st := SwapExactTokensForTokensParams{
		AmountIn:     amountIn,
		AmountOutMin: amountOutMin,
		Path:         path,
		To:           to,
		Deadline:     deadline,
	}
	pair, err := u.PairToPrices(context.Background(), path)
	if err != nil {
		return
	}
	initialPair := pair
	tf := st.BinarySearch(pair)
	tf.InitialPair = initialPair
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapExactTokensForTokens==================================")
		u.PrintTradeSummaries(tx, tf, pair, path[0].String(), st.AmountIn, st.AmountOutMin)
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount.String(), "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit.String())
		fmt.Println("sandwich: ====================================SwapExactTokensForTokens==================================")
	}
	u.SwapExactTokensForTokensParamsSlice = append(u.SwapExactTokensForTokensParamsSlice, st)
}

func (u *UniswapV2Client) PairToPrices(ctx context.Context, pairAddr []common.Address) (UniswapV2Pair, error) {
	if len(pairAddr) == 2 {
		pairContractAddr := u.GetPairContractFromFactory(ctx, pairAddr[0].String(), pairAddr[1].String())
		return u.GetPairContractPrices(ctx, pairContractAddr.String())
	}
	return UniswapV2Pair{}, errors.New("pair address length is not 2")
}

func (u *UniswapV2Client) SwapTokensForExactTokens(tx MevTx, args map[string]interface{}) {
	amountOut, err := ParseBigInt(args["amountOut"])
	if err != nil {
		return
	}
	amountInMax, err := ParseBigInt(args["amountInMax"])
	if err != nil {
		return
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		return
	}
	to, err := ConvertToAddress(args["to"])
	if err != nil {
		return
	}
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		return
	}
	st := SwapTokensForExactTokensParams{
		AmountOut:   amountOut,
		AmountInMax: amountInMax,
		Path:        path,
		To:          to,
		Deadline:    deadline,
	}
	pair, err := u.PairToPrices(context.Background(), path)
	if err != nil {
		return
	}
	initialPair := pair
	tf := st.BinarySearch(pair)
	tf.InitialPair = initialPair
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapTokensForExactTokens==================================")
		u.PrintTradeSummaries(tx, tf, pair, path[0].String(), st.AmountInMax, st.AmountOut)
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount.String(), "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit.String())
		fmt.Println("sandwich: ====================================SwapTokensForExactTokens==================================")
	}
	u.SwapTokensForExactTokensParamsSlice = append(u.SwapTokensForExactTokensParamsSlice, st)
}

func (u *UniswapV2Client) SwapExactETHForTokens(tx MevTx, args map[string]interface{}, payableEth *big.Int) {
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		return
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		return
	}
	to, err := ConvertToAddress(args["to"])
	if err != nil {
		return
	}
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		return
	}
	st := SwapExactETHForTokensParams{
		AmountOutMin: amountOutMin,
		Path:         path,
		To:           to,
		Deadline:     deadline,
		Value:        payableEth,
	}

	pair, err := u.PairToPrices(context.Background(), path)
	if err != nil {
		return
	}
	initialPair := pair
	tf := st.BinarySearch(pair)
	tf.InitialPair = initialPair
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapExactETHForTokens==================================")
		u.PrintTradeSummaries(tx, tf, pair, path[0].String(), st.Value, st.AmountOutMin)
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount.String(), "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit.String())
		fmt.Println("sandwich: ====================================SwapExactETHForTokens==================================")
	}
	u.SwapExactETHForTokensParamsSlice = append(u.SwapExactETHForTokensParamsSlice, st)
}

func (u *UniswapV2Client) SwapTokensForExactETH(tx MevTx, args map[string]interface{}) {
	amountOut, err := ParseBigInt(args["amountOut"])
	if err != nil {
		return
	}
	amountInMax, err := ParseBigInt(args["amountInMax"])
	if err != nil {
		return
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		return
	}
	to, err := ConvertToAddress(args["to"])
	if err != nil {
		return
	}
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		return
	}
	st := SwapTokensForExactETHParams{
		AmountOut:   amountOut,
		AmountInMax: amountInMax,
		Path:        path,
		To:          to,
		Deadline:    deadline,
	}
	pair, err := u.PairToPrices(context.Background(), path)
	if err != nil {
		return
	}
	initialPair := pair
	tf := st.BinarySearch(pair)
	tf.InitialPair = initialPair
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapTokensForExactETH==================================")
		u.PrintTradeSummaries(tx, tf, pair, path[0].String(), st.AmountInMax, st.AmountOut)
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount.String(), "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit.String())
		fmt.Println("sandwich: ====================================SwapTokensForExactETH==================================")
	}
	u.SwapTokensForExactETHParamsSlice = append(u.SwapTokensForExactETHParamsSlice, st)
}

func (u *UniswapV2Client) SwapExactTokensForETH(tx MevTx, args map[string]interface{}) {
	amountIn, err := ParseBigInt(args["amountIn"])
	if err != nil {
		return
	}
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		return
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		return
	}
	to, err := ConvertToAddress(args["to"])
	if err != nil {
		return
	}
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		return
	}
	st := SwapExactTokensForETHParams{
		AmountIn:     amountIn,
		AmountOutMin: amountOutMin,
		Path:         path,
		To:           to,
		Deadline:     deadline,
	}
	pair, err := u.PairToPrices(context.Background(), path)
	if err != nil {
		return
	}
	initialPair := pair
	tf := st.BinarySearch(pair)
	tf.InitialPair = initialPair
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapExactTokensForETH==================================")
		u.PrintTradeSummaries(tx, tf, pair, path[0].String(), st.AmountIn, st.AmountOutMin)
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount.String(), "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit.String())
		fmt.Println("sandwich: ====================================SwapExactTokensForETH==================================")
	}
	u.SwapExactTokensForETHParamsSlice = append(u.SwapExactTokensForETHParamsSlice, st)
}

func (u *UniswapV2Client) SwapETHForExactTokens(tx MevTx, args map[string]interface{}, payableEth *big.Int) {
	amountOut, err := ParseBigInt(args["amountOut"])
	if err != nil {
		return
	}
	path, err := ConvertToAddressSlice(args["path"])
	if err != nil {
		return
	}
	to, err := ConvertToAddress(args["to"])
	if err != nil {
		return
	}
	deadline, err := ParseBigInt(args["deadline"])
	if err != nil {
		return
	}
	st := SwapETHForExactTokensParams{
		AmountOut: amountOut,
		Path:      path,
		To:        to,
		Deadline:  deadline,
		Value:     payableEth,
	}
	pair, err := u.PairToPrices(context.Background(), path)
	if err != nil {
		return
	}
	initialPair := pair
	tf := st.BinarySearch(pair)
	tf.InitialPair = initialPair
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapETHForExactTokens==================================")
		u.PrintTradeSummaries(tx, tf, pair, path[0].String(), st.Value, st.AmountOut)
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount.String(), "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit.String())
		fmt.Println("sandwich: ====================================SwapETHForExactTokens==================================")
	}
	u.SwapETHForExactTokensParamsSlice = append(u.SwapETHForExactTokensParamsSlice, st)
}
