package web3_client

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
)

type ExactOutputParams struct {
	Path            []byte           `json:"path"`
	Recipient       accounts.Address `json:"recipient"`
	AmountOut       *big.Int         `json:"amountOut"`
	AmountInMaximum *big.Int         `json:"amountInMaximum"`

	TokenFeePath TokenFeePath `json:"tokenFeePath,omitempty"`
}

type JSONExactOutputParams struct {
	Path            []byte           `json:"path"`
	Recipient       accounts.Address `json:"recipient"`
	AmountOut       string           `json:"amountOut"`
	AmountInMaximum string           `json:"amountInMaximum"`

	TokenFeePath TokenFeePath `json:"tokenFeePath,omitempty"`
}

func (o *ExactOutputParams) BinarySearch(pd *PricingData) TradeExecutionFlowJSON {
	low := big.NewInt(0)
	high := new(big.Int).Set(o.AmountInMaximum)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlowJSON{
		Trade: Trade{
			TradeMethod:           exactOutput,
			JSONExactOutputParams: o.ConvertToJSONType(),
		},
	}
	frontRunTokenIn := pd.V3Pair.Token0
	sandwichTokenIn := pd.V3Pair.Token1
	if o.TokenFeePath.TokenIn.Hex() == pd.V3Pair.Token1.Address.Hex() {
		frontRunTokenIn = pd.V3Pair.Token1
		sandwichTokenIn = pd.V3Pair.Token0
	}
	for low.Cmp(high) <= 0 {
		mockPairResp := pd.V3Pair
		mid = new(big.Int).Add(low, high)
		mid = DivideByHalf(mid)
		// Front run trade
		amountInFrontRun := mid
		toFrontRun, _, err := mockPairResp.PriceImpact(ctx, frontRunTokenIn, amountInFrontRun)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		// User trade
		userTrade, _, err := mockPairResp.PriceImpact(ctx, frontRunTokenIn, o.AmountInMaximum)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		difference := new(big.Int).Sub(userTrade.Quotient(), o.AmountOut)
		// if diff <= 0 then it searches left
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(mid, big.NewInt(1))
			continue
		}
		// Sandwich trade
		toSandwich, _, err := mockPairResp.PriceImpact(ctx, sandwichTokenIn, toFrontRun.Quotient())
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf
		}
		profit := new(big.Int).Sub(toSandwich.Quotient(), toFrontRun.Quotient())
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			tokenSellAmountAtMaxProfit = mid
			tf.FrontRunTrade = artemis_trading_types.JSONTradeOutcome{
				AmountIn:      amountInFrontRun.String(),
				AmountInAddr:  frontRunTokenIn.Address,
				AmountOut:     toFrontRun.Quotient().String(),
				AmountOutAddr: sandwichTokenIn.Address,
			}
			tf.UserTrade = artemis_trading_types.JSONTradeOutcome{
				AmountIn:      o.AmountInMaximum.String(),
				AmountInAddr:  frontRunTokenIn.Address,
				AmountOut:     userTrade.Quotient().String(),
				AmountOutAddr: sandwichTokenIn.Address,
			}
			tf.SandwichTrade = artemis_trading_types.JSONTradeOutcome{
				AmountIn:      toFrontRun.Quotient().String(),
				AmountInAddr:  sandwichTokenIn.Address,
				AmountOut:     toSandwich.Quotient().String(),
				AmountOutAddr: frontRunTokenIn.Address,
			}
		}
		// If profit is negative, reduce the high boundary
		if profit.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(mid, big.NewInt(1))
		} else {
			// If profit is positive, increase the low boundary
			low = new(big.Int).Add(mid, big.NewInt(1))
		}
	}
	sp := SandwichTradePrediction{
		SellAmount:     tokenSellAmountAtMaxProfit,
		ExpectedProfit: maxProfit,
	}
	tf.SandwichPrediction = sp.ConvertToJSONType()
	return tf
}

func (o *ExactOutputParams) Decode(ctx context.Context, args map[string]interface{}) error {
	params, ok := args["params"].(struct {
		Path            []byte         "json:\"path\""
		Recipient       common.Address "json:\"recipient\""
		AmountOut       *big.Int       "json:\"amountOut\""
		AmountInMaximum *big.Int       "json:\"amountInMaximum\""
	})
	if !ok {
		return errors.New("invalid params")
	}
	hexStr := accounts.Bytes2Hex(params.Path)
	tfp := TokenFeePath{
		TokenIn: accounts.HexToAddress(hexStr[:40]),
	}
	var pathList []TokenFee
	for i := 0; i < len(hexStr[40:]); i += 46 {
		fee, _ := new(big.Int).SetString(hexStr[40:][i:i+6], 16)
		token := accounts.HexToAddress(hexStr[40:][i+6 : i+46])
		tf := TokenFee{
			Token: token,
			Fee:   fee,
		}
		pathList = append(pathList, tf)
	}
	tfp.Path = pathList
	o.TokenFeePath = tfp
	o.Path = params.Path
	o.Recipient = accounts.Address(params.Recipient)
	o.AmountOut = params.AmountOut
	o.AmountInMaximum = params.AmountInMaximum
	return nil
}

func (o *JSONExactOutputParams) ConvertToBigIntType() *ExactOutputParams {
	amountInMax, _ := new(big.Int).SetString(o.AmountInMaximum, 10)
	amountOut, _ := new(big.Int).SetString(o.AmountOut, 10)
	return &ExactOutputParams{
		AmountInMaximum: amountInMax,
		AmountOut:       amountOut,
		Path:            o.Path,
		Recipient:       o.Recipient,
		TokenFeePath:    o.TokenFeePath,
	}
}

func (o *ExactOutputParams) ConvertToJSONType() *JSONExactOutputParams {
	return &JSONExactOutputParams{
		AmountInMaximum: o.AmountInMaximum.String(),
		AmountOut:       o.AmountOut.String(),
		Path:            o.Path,
		Recipient:       o.Recipient,
		TokenFeePath:    o.TokenFeePath,
	}
}
