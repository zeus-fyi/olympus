package web3_client

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type TradeSummary struct {
	Tx            MevTx
	Pd            *uniswap_pricing.UniswapPricingData
	Tf            TradeExecutionFlowJSON
	TokenAddr     string
	BuyWithAmount *big.Int
	MinimumAmount *big.Int
}

func (u *UniswapClient) PrintTradeSummaries(ts *TradeSummary) {
	log.Info().Msgf("TradeSummary")
	jsonTX := artemis_trading_types.JSONTx{}
	err := jsonTX.UnmarshalTx(ts.Tx.Tx)
	if err != nil {
		panic(err)
	}
	ts.Tf.Tx = jsonTX
	wc := web3_actions.NewWeb3ActionsClient(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.NodeURL)
	wc.Dial()
	bn, ber := wc.C.BlockNumber(ctx)
	if ber != nil {
		log.Err(ber)
	}
	fmt.Println("blockNumber:", bn)
	defer wc.Close()
	pair := ts.Pd.V2Pair
	ts.Tf.CurrentBlockNumber = new(big.Int).SetUint64(bn)
	expectedOut, err := pair.GetQuoteUsingTokenAddr(ts.TokenAddr, ts.BuyWithAmount)
	if err != nil {
		fmt.Println("GetQuoteUsingTokenAddr", err)
		return
	}
	diff := new(big.Int).Sub(expectedOut, ts.MinimumAmount)
	purchasedTokenAddr := pair.GetOppositeToken(ts.TokenAddr).String()
	if u.PrintDetails {
		fmt.Printf("Token0 Address: %s Token0 Reserve: %s,\nToken1 Address %s, Token1 Reserve: %s\n", pair.Token0.String(), pair.Reserve0.String(), pair.Token1.String(), pair.Reserve1.String())
		fmt.Printf("Expected amount %s %s token from trade at current rate \n", expectedOut.String(), purchasedTokenAddr)
		fmt.Printf("BuyWithAmount minimum %s %s token needed from trade \n", ts.MinimumAmount.String(), purchasedTokenAddr)
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
		sender := types.LatestSignerForChainID(ts.Tx.Tx.ChainId())
		from, ferr := sender.Sender(ts.Tx.Tx)
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
		err = artemis_mev_models.InsertMempoolTx(ctx, txMempool)
		if err != nil {
			fmt.Printf("InsertMempoolTx err: %s", err)
			return
		}
	} else {
		if u.PrintDetails {
			fmt.Printf("Negative difference between expected and minimum amount is %s %s token \n", diff.String(), ts.TokenAddr)
		}
	}
	if ts.MinimumAmount.Cmp(big.NewInt(0)) == 0 {
		fmt.Printf("BuyWithAmount minimum is 0, so no trade will be executed \n")
		return
	}
	if u.PrintDetails {
		slippage := new(big.Int).Mul(diff, big.NewInt(100))
		slippagePercent := new(big.Int).Div(slippage, ts.MinimumAmount)
		fmt.Printf("Slippage is %s %% \n", slippagePercent.String())
		fmt.Printf("Buy %s %s token for %s %s token \n\n", expectedOut.String(), pair.GetOppositeToken(ts.TokenAddr).String(), ts.BuyWithAmount.String(), ts.TokenAddr)
	}
	return
}

//
//func (u *UniswapClient) PrintTradeSummaries(tx MevTx, tf TradeExecutionFlowJSON, pair UniswapV2Pair, tokenAddr string, amount, amountMin *big.Int) {
//	tf.Tx = tx.Tx
//	u.Web3Client.Dial()
//	defer u.Web3Client.Close()
//	bn, err := u.Web3Client.GetHeadBlockHeight(ctx)
//	if err != nil {
//		fmt.Println("GetBlockNumber Error", err)
//		return
//	}
//	tf.CurrentBlockNumber = bn
//	expectedOut, err := pair.GetQuoteUsingTokenAddr(tokenAddr, amount)
//	if err != nil {
//		fmt.Println("GetQuoteUsingTokenAddr", err)
//		return
//	}
//	diff := new(big.Int).Sub(expectedOut, amountMin)
//	purchasedTokenAddr := pair.GetOppositeToken(tokenAddr).String()
//	if u.PrintDetails {
//		fmt.Printf("Token0 Address: %s Token0 Reserve: %s,\nToken1 Address %s, Token1 Reserve: %s\n", pair.Token0.String(), pair.Reserve0.String(), pair.Token1.String(), pair.Reserve1.String())
//		fmt.Printf("Expected amount %s %s token from trade at current rate \n", expectedOut.String(), purchasedTokenAddr)
//		fmt.Printf("BuyWithAmount minimum %s %s token needed from trade \n", amountMin.String(), purchasedTokenAddr)
//	}
//
//	if u.BlockNumber.String() != tf.CurrentBlockNumber.String() {
//		log.Info().Interface("currentBlockNumber", tf.CurrentBlockNumber.String()).Interface("startingBlockNumber", u.BlockNumber.String()).Msg("block number transition exiting due to stale data")
//		return
//	}
//	if diff.Cmp(big.NewInt(0)) == 1 {
//		fmt.Printf("Positive difference between expected and minimum amount is %s %s token \n", diff.String(), tokenAddr)
//		b, berr := json.MarshalIndent(tf, "", "  ")
//		if berr != nil {
//			return
//		}
//		if u.PrintLocal {
//			u.Path.FnOut = fmt.Sprintf("%s-%d.json", tf.Trade.TradeMethod, u.BlockNumber)
//			err = u.Path.WriteToFileOutPath(b)
//			if err != nil {
//				return
//			}
//		}
//
//		btf, berr := json.Marshal(tf)
//		if berr != nil {
//			return
//		}
//		b, berr = json.Marshal(tf.Tx)
//		if berr != nil {
//			return
//		}
//		fromStr := ""
//		sender := types.LatestSignerForChainID(tx.Tx.ChainId())
//		from, ferr := sender.Sender(tx.Tx)
//		if ferr != nil {
//			log.Err(err).Msg("failed to get sender")
//		} else {
//			fromStr = from.String()
//		}
//
//		txMempool := artemis_autogen_bases.EthMempoolMevTx{
//			ProtocolNetworkID: hestia_req_types.EthereumMainnetProtocolNetworkID,
//			Tx:                string(b),
//			TxFlowPrediction:  string(btf),
//			TxHash:            tx.Tx.Hash().String(),
//			Nonce:             int(tx.Tx.Nonce()),
//			From:              fromStr,
//			To:                tx.Tx.To().String(),
//			BlockNumber:       int(u.BlockNumber.Int64()),
//		}
//		u.Trades = append(u.Trades, txMempool)
//		err = artemis_validator_service_groups_models.InsertMempoolTx(ctx, txMempool)
//		if err != nil {
//			fmt.Printf("InsertMempoolTx err: %s", err)
//			return
//		}
//	} else {
//		if u.PrintDetails {
//			fmt.Printf("Negative difference between expected and minimum amount is %s %s token \n", diff.String(), tokenAddr)
//		}
//	}
//	if amountMin.Cmp(big.NewInt(0)) == 0 {
//		fmt.Printf("BuyWithAmount minimum is 0, so no trade will be executed \n")
//		return
//	}
//	if u.PrintDetails {
//		slippage := new(big.Int).Mul(diff, big.NewInt(100))
//		slippagePercent := new(big.Int).Div(slippage, amountMin)
//		fmt.Printf("Slippage is %s %% \n", slippagePercent.String())
//		fmt.Printf("Buy %s %s token for %s %s token \n\n", expectedOut.String(), pair.GetOppositeToken(tokenAddr).String(), amount.String(), tokenAddr)
//	}
//	return
//}
