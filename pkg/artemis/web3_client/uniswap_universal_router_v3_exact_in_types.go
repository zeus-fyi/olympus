package web3_client

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
)

const (
	V3SwapExactIn = "V3_SWAP_EXACT_IN"
)

type V3SwapExactInParams struct {
	AmountIn     *big.Int                           `json:"amountIn"`
	AmountOutMin *big.Int                           `json:"amountOutMin"`
	Path         artemis_trading_types.TokenFeePath `json:"path"`
	To           accounts.Address                   `json:"to"`
	PayerIsUser  bool                               `json:"payerIsUser"`
}

type JSONV3SwapExactInParams struct {
	AmountIn     string                             `json:"amountIn"`
	AmountOutMin string                             `json:"amountOutMin"`
	Path         artemis_trading_types.TokenFeePath `json:"path"`
	To           accounts.Address                   `json:"to"`
	PayerIsUser  bool                               `json:"payerIsUser"`
}

func (s *V3SwapExactInParams) Encode(ctx context.Context, abiFile *abi.ABI) ([]byte, error) {
	if abiFile == nil {
		inputs, err := UniversalRouterDecoderAbi.Methods[V3SwapExactIn].Inputs.Pack(s.To, s.AmountIn, s.AmountOutMin, s.Path.Encode(), s.PayerIsUser)
		if err != nil {
			return nil, err
		}
		return inputs, nil
	} else {
		inputs, err := abiFile.Methods[V3SwapExactIn].Inputs.Pack(s.To, s.AmountIn, s.AmountOutMin, s.Path.Encode(), s.PayerIsUser)
		if err != nil {
			return nil, err
		}
		return inputs, nil
	}
}

func (s *V3SwapExactInParams) Decode(ctx context.Context, data []byte, abiFile *abi.ABI) error {
	args := make(map[string]interface{})
	if abiFile == nil {
		err := UniversalRouterDecoderAbi.Methods[V3SwapExactIn].Inputs.UnpackIntoMap(args, data)
		if err != nil {
			log.Err(err).Msg("V3SwapExactInParams: UniversalRouterDecoderAbi Decode failed to unpack")
			return err
		}
	} else {
		err := abiFile.Methods[V3SwapExactIn].Inputs.UnpackIntoMap(args, data)
		if err != nil {
			log.Err(err).Msg("V3SwapExactInParams: abiFile Decode failed to unpack")
			return err
		}
	}
	amountIn, err := ParseBigInt(args["amountIn"])
	if err != nil {
		log.Warn().Msg("V3SwapExactInParams: Decode failed to parse amountIn")
		return err
	}
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		log.Warn().Msg("V3SwapExactInParams: Decode failed to parse amountOutMin")
		return err
	}
	pathBytes := args["path"].([]byte)
	hexStr := accounts.Bytes2Hex(pathBytes)
	tfp := artemis_trading_types.TokenFeePath{
		TokenIn: accounts.HexToAddress(hexStr[:40]),
	}
	var pathList []artemis_trading_types.TokenFee
	for i := 0; i < len(hexStr[40:]); i += 46 {
		fee, ok := new(big.Int).SetString(hexStr[40:][i:i+6], 16)
		if !ok {
			log.Warn().Msg("V3SwapExactInParams: Decode failed to parse fee")
			return errors.New("V3SwapExactInParams: Decode failed to parse fee")
		}
		token := accounts.HexToAddress(hexStr[40:][i+6 : i+46])
		tf := artemis_trading_types.TokenFee{
			Token: token,
			Fee:   fee,
		}
		pathList = append(pathList, tf)
	}
	tfp.Path = pathList

	to, err := ConvertToAddress(args["recipient"])
	if err != nil {
		return err
	}
	payerIsSender := args["payerIsUser"].(bool)
	s.AmountIn = amountIn
	s.AmountOutMin = amountOutMin
	s.Path = tfp
	s.To = to
	s.PayerIsUser = payerIsSender
	return err
}

func (s *JSONV3SwapExactInParams) ConvertToBigIntType() *V3SwapExactInParams {
	amountIn, _ := new(big.Int).SetString(s.AmountIn, 10)
	amountOutMin, _ := new(big.Int).SetString(s.AmountOutMin, 10)
	return &V3SwapExactInParams{
		AmountIn:     amountIn,
		AmountOutMin: amountOutMin,
		Path:         s.Path,
		To:           s.To,
		PayerIsUser:  s.PayerIsUser,
	}
}

func (s *V3SwapExactInParams) ConvertToJSONType() *JSONV3SwapExactInParams {
	return &JSONV3SwapExactInParams{
		AmountIn:     s.AmountIn.String(),
		AmountOutMin: s.AmountOutMin.String(),
		Path:         s.Path,
		To:           s.To,
		PayerIsUser:  s.PayerIsUser,
	}
}
