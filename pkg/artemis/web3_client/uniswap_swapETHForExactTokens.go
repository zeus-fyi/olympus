package web3_client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (u *UniswapClient) SwapETHForExactTokens(tx MevTx, args map[string]interface{}, payableEth *big.Int) {
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
	pd, err := u.GetPricingData(ctx, path)
	if err != nil {
		return
	}
	initialPair := pd.v2Pair
	tf := st.BinarySearch(pd.v2Pair)
	tf.InitialPair = initialPair.ConvertToJSONType()
	if u.PrintOn {
		fmt.Println("\nsandwich: ==================================SwapETHForExactTokens==================================")
		ts := TradeSummary{
			Tx:        tx,
			Pd:        pd,
			Tf:        tf,
			TokenAddr: path[0].String(),
			Amount:    st.Value,
			AmountMin: st.AmountOut,
		}
		u.PrintTradeSummaries2(&ts)
		fmt.Println("Sell Token: ", path[0].String(), "Buy Token", path[1].String(), "Sell Amount: ", tf.SandwichPrediction.SellAmount, "Expected Profit: ", tf.SandwichPrediction.ExpectedProfit)
		fmt.Println("sandwich: ====================================SwapETHForExactTokens==================================")
	}
	u.SwapETHForExactTokensParamsSlice = append(u.SwapETHForExactTokensParamsSlice, st)
}

type TradeSummary struct {
	Tx        MevTx
	Pd        *PricingData
	Tf        TradeExecutionFlowJSON
	TokenAddr string
	Amount    *big.Int
	AmountMin *big.Int
}

func (u *UniswapClient) PrintTradeSummaries2(ts *TradeSummary) {
	ts.Tf.Tx = ts.Tx.Tx
	u.Web3Client.Dial()
	defer u.Web3Client.Close()
	bn, err := u.Web3Client.GetHeadBlockHeight(ctx)
	if err != nil {
		fmt.Println("GetBlockNumber Error", err)
		return
	}
	pair := ts.Pd.v2Pair
	ts.Tf.CurrentBlockNumber = bn
	expectedOut, err := pair.GetQuoteUsingTokenAddr(ts.TokenAddr, ts.Amount)
	if err != nil {
		fmt.Println("GetQuoteUsingTokenAddr", err)
		return
	}
	diff := new(big.Int).Sub(expectedOut, ts.AmountMin)
	purchasedTokenAddr := pair.GetOppositeToken(ts.TokenAddr).String()
	if u.PrintDetails {
		fmt.Printf("Token0 Address: %s Token0 Reserve: %s,\nToken1 Address %s, Token1 Reserve: %s\n", pair.Token0.String(), pair.Reserve0.String(), pair.Token1.String(), pair.Reserve1.String())
		fmt.Printf("Expected amount %s %s token from trade at current rate \n", expectedOut.String(), purchasedTokenAddr)
		fmt.Printf("Amount minimum %s %s token needed from trade \n", ts.AmountMin.String(), purchasedTokenAddr)
	}

	if u.BlockNumber.String() != ts.Tf.CurrentBlockNumber.String() {
		log.Info().Interface("currentBlockNumber", ts.Tf.CurrentBlockNumber.String()).Interface("startingBlockNumber", u.BlockNumber.String()).Msg("block number transition exiting due to stale data")
		return
	}
	if diff.Cmp(big.NewInt(0)) == 1 {
		fmt.Printf("Positive difference between expected and minimum amount is %s %s token \n", diff.String(), ts.TokenAddr)
		b, berr := json.MarshalIndent(ts.Tf, "", "  ")
		if berr != nil {
			return
		}
		if u.PrintLocal {
			u.Path.FnOut = fmt.Sprintf("%s-%d.json", ts.Tf.Trade.TradeMethod, u.BlockNumber)
			err = u.Path.WriteToFileOutPath(b)
			if err != nil {
				return
			}
		}

		btf, berr := json.Marshal(ts.Tf)
		if berr != nil {
			return
		}
		b, berr = json.Marshal(ts.Tf.Tx)
		if berr != nil {
			return
		}
		fromStr := ""
		sender := types.LatestSignerForChainID(ts.Tf.Tx.ChainId())
		from, ferr := sender.Sender(ts.Tf.Tx)
		if ferr != nil {
			log.Err(err).Msg("failed to get sender")
		} else {
			fromStr = from.String()
		}

		txMempool := artemis_autogen_bases.EthMempoolMevTx{
			ProtocolNetworkID: hestia_req_types.EthereumMainnetProtocolNetworkID,
			Tx:                string(b),
			TxFlowPrediction:  string(btf),
			TxHash:            ts.Tx.Tx.Hash().String(),
			Nonce:             int(ts.Tx.Tx.Nonce()),
			From:              fromStr,
			To:                ts.Tx.Tx.To().String(),
			BlockNumber:       int(u.BlockNumber.Int64()),
		}
		u.Trades = append(u.Trades, txMempool)
		err = artemis_validator_service_groups_models.InsertMempoolTx(ctx, txMempool)
		if err != nil {
			fmt.Printf("InsertMempoolTx err: %s", err)
			return
		}
	} else {
		if u.PrintDetails {
			fmt.Printf("Negative difference between expected and minimum amount is %s %s token \n", diff.String(), ts.TokenAddr)
		}
	}
	if ts.AmountMin.Cmp(big.NewInt(0)) == 0 {
		fmt.Printf("Amount minimum is 0, so no trade will be executed \n")
		return
	}
	if u.PrintDetails {
		slippage := new(big.Int).Mul(diff, big.NewInt(100))
		slippagePercent := new(big.Int).Div(slippage, ts.AmountMin)
		fmt.Printf("Slippage is %s %% \n", slippagePercent.String())
		fmt.Printf("Buy %s %s token for %s %s token \n\n", expectedOut.String(), pair.GetOppositeToken(ts.TokenAddr).String(), ts.Amount.String(), ts.TokenAddr)
	}
	return
}

type PricingData struct {
	v2Pair UniswapV2Pair
}

func (u *UniswapClient) GetPricingData(ctx context.Context, path []accounts.Address) (*PricingData, error) {
	pair, err := u.PairToPrices(ctx, path)
	if err != nil {
		return nil, err
	}
	return &PricingData{
		v2Pair: pair,
	}, nil
}
