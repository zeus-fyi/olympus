package web3_client

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	uniswap_pricing "github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/uniswap"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
)

type ExactInputParams struct {
	Path             []byte           `json:"path"`
	Recipient        accounts.Address `json:"recipient"`
	AmountIn         *big.Int         `json:"amountIn"`
	AmountOutMinimum *big.Int         `json:"amountOutMinimum"`

	TokenFeePath artemis_trading_types.TokenFeePath `json:"tokenFeePath,omitempty"`
}
type JSONExactInputParams struct {
	Path             []byte           `json:"path"`
	Recipient        accounts.Address `json:"recipient"`
	AmountIn         string           `json:"amountIn"`
	AmountOutMinimum string           `json:"amountOutMinimum"`

	TokenFeePath artemis_trading_types.TokenFeePath `json:"tokenFeePath,omitempty"`
}

func (in *JSONExactInputParams) ConvertToBigIntType() (*ExactInputParams, error) {
	amountIn, ok := new(big.Int).SetString(in.AmountIn, 10)
	if !ok {
		log.Warn().Msg("JSONExactInputParams: could not convert amountIn to big.Int")
		return nil, errors.New("could not convert amountIn to big.Int")
	}
	amountOutMin, ok := new(big.Int).SetString(in.AmountOutMinimum, 10)
	if !ok {
		log.Warn().Msg("JSONExactInputParams: could not convert amountOutMinimum to big.Int")
		return nil, errors.New("could not convert amountOutMinimum to big.Int")
	}
	return &ExactInputParams{
		AmountIn:         amountIn,
		AmountOutMinimum: amountOutMin,
		Path:             in.Path,
		Recipient:        in.Recipient,
		TokenFeePath:     in.TokenFeePath,
	}, nil
}

func (in *ExactInputParams) ConvertToJSONType() *JSONExactInputParams {
	return &JSONExactInputParams{
		AmountIn:         in.AmountIn.String(),
		AmountOutMinimum: in.AmountOutMinimum.String(),
		Path:             in.Path,
		Recipient:        in.Recipient,
		TokenFeePath:     in.TokenFeePath,
	}
}

func (in *ExactInputParams) BinarySearch(pd *uniswap_pricing.UniswapPricingData) (TradeExecutionFlow, error) {
	low := big.NewInt(0)
	high := new(big.Int).Set(in.AmountIn)
	var mid *big.Int
	var maxProfit *big.Int
	var tokenSellAmountAtMaxProfit *big.Int
	tf := TradeExecutionFlow{
		InitialPairV3: &pd.V3Pair,
		Trade: Trade{
			TradeMethod:          exactInput,
			JSONExactInputParams: in.ConvertToJSONType(),
		},
	}
	frontRunTokenIn := pd.V3Pair.Token0
	sandwichTokenIn := pd.V3Pair.Token1
	if in.TokenFeePath.TokenIn.Hex() == pd.V3Pair.Token1.Address.Hex() {
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
			return tf, err
		}
		// User trade
		userTrade, _, err := mockPairResp.PriceImpact(ctx, frontRunTokenIn, in.AmountIn)
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf, err
		}
		difference := new(big.Int).Sub(userTrade.Quotient(), in.AmountOutMinimum)
		if difference.Cmp(big.NewInt(0)) < 0 {
			high = new(big.Int).Sub(mid, big.NewInt(1))
			continue
		}
		// Sandwich trade
		toSandwich, _, err := mockPairResp.PriceImpact(ctx, sandwichTokenIn, toFrontRun.Quotient())
		if err != nil {
			log.Err(err).Msg("error in price impact")
			return tf, err
		}
		profit := new(big.Int).Sub(toSandwich.Quotient(), toFrontRun.Quotient())
		if maxProfit == nil || profit.Cmp(maxProfit) > 0 {
			maxProfit = profit
			tokenSellAmountAtMaxProfit = mid
			tf.FrontRunTrade = artemis_trading_types.TradeOutcome{
				AmountIn:      amountInFrontRun,
				AmountInAddr:  frontRunTokenIn.Address,
				AmountOut:     toFrontRun.Quotient(),
				AmountOutAddr: sandwichTokenIn.Address,
			}
			tf.UserTrade = artemis_trading_types.TradeOutcome{
				AmountIn:      in.AmountIn,
				AmountInAddr:  frontRunTokenIn.Address,
				AmountOut:     userTrade.Quotient(),
				AmountOutAddr: sandwichTokenIn.Address,
			}
			tf.SandwichTrade = artemis_trading_types.TradeOutcome{
				AmountIn:      toFrontRun.Quotient(),
				AmountInAddr:  sandwichTokenIn.Address,
				AmountOut:     toSandwich.Quotient(),
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
	tf.SandwichPrediction = sp
	return tf, nil
}

func (in *ExactInputParams) Decode(ctx context.Context, args map[string]interface{}) error {
	params, ok := args["params"].(struct {
		Path             []byte         "json:\"path\""
		Recipient        common.Address "json:\"recipient\""
		AmountIn         *big.Int       "json:\"amountIn\""
		AmountOutMinimum *big.Int       "json:\"amountOutMinimum\""
	})
	if !ok {
		log.Warn().Msg("ExactInputParams: Decode: invalid params")
		return errors.New("invalid params")
	}
	hexStr := accounts.Bytes2Hex(params.Path)
	tfp := artemis_trading_types.TokenFeePath{
		TokenIn: accounts.HexToAddress(hexStr[:40]),
	}
	var pathList []artemis_trading_types.TokenFee
	for i := 0; i < len(hexStr[40:]); i += 46 {
		fee, fok := new(big.Int).SetString(hexStr[40:][i:i+6], 16)
		if !fok {
			log.Warn().Msg("ExactInputParams: Decode: invalid fee")
			return errors.New("ExactInputParams: invalid fee")
		}
		token := accounts.HexToAddress(hexStr[40:][i+6 : i+46])
		tf := artemis_trading_types.TokenFee{
			Token: token,
			Fee:   fee,
		}
		pathList = append(pathList, tf)
	}
	tfp.Path = pathList
	in.TokenFeePath = tfp
	in.Path = params.Path
	in.Recipient = accounts.Address(params.Recipient)
	in.AmountIn = params.AmountIn
	in.AmountOutMinimum = params.AmountOutMinimum
	return nil
}
