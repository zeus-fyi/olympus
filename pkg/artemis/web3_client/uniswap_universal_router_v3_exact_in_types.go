package web3_client

import (
	"context"
	"errors"
	"fmt"
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
	inputs, err := UniversalRouterDecoderAbi.Methods[V3SwapExactIn].Inputs.Pack(s.To, s.AmountIn, s.AmountOutMin, s.Path.Encode(), s.PayerIsUser)
	if err != nil {
		log.Warn().Err(err).Msg("V3SwapExactInParams: UniversalRouterDecoderAbi Encode failed to pack")
		return nil, err
	}
	return inputs, nil
}

func (s *V3SwapExactInParams) Decode(ctx context.Context, data []byte, abiFile *abi.ABI) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[V3SwapExactIn].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		log.Err(err).Msg("V3SwapExactInParams: UniversalRouterDecoderAbi Decode failed to unpack")
		return err
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
	pathInterface, pok := args["path"]
	if !pok || pathInterface == nil {
		// Handle the situation when args["path"] doesn't exist or is nil
		log.Warn().Msg("V3SwapExactInParams: 'path' does not exist or is nil")
		return fmt.Errorf("path does not exist or is nil")
	}
	pathBytes, ok := pathInterface.([]byte)
	if !ok {
		// Handle the situation when the conversion fails
		log.Warn().Msg("V3SwapExactInParams: failed to convert 'path' to []byte")
		return fmt.Errorf("failed to convert path to []byte")
	}
	hexStr := accounts.Bytes2Hex(pathBytes)
	tfp := artemis_trading_types.TokenFeePath{
		TokenIn: accounts.HexToAddress(hexStr[:40]),
	}
	var pathList []artemis_trading_types.TokenFee
	for i := 0; i < len(hexStr[40:]); i += 46 {
		fee, fok := new(big.Int).SetString(hexStr[40:][i:i+6], 16)
		if !fok {
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
		log.Err(err).Msg("V3SwapExactInParams: Decode failed to parse recipient")
		return err
	}
	payerIsUserInterface, pok := args["payerIsUser"]
	if !pok || payerIsUserInterface == nil {
		// Handle the situation when args["path"] doesn't exist or is nil
		log.Warn().Msg("V3SwapExactInParams: 'payerIsUser' does not exist or is nil defaulting to false")
		s.PayerIsUser = false
	} else {
		payerIsUserBool, ok1 := payerIsUserInterface.(bool)
		if !ok1 {
			// Handle the situation when the conversion fails
			log.Warn().Msg("V3SwapExactInParams: failed to convert 'payerIsUser' to bool")
			log.Warn().Msg("V3SwapExactInParams: 'payerIsUser' does not exist or is nil defaulting to false")
		}
		s.PayerIsUser = payerIsUserBool
	}

	s.AmountIn = amountIn
	s.AmountOutMin = amountOutMin
	s.Path = tfp
	s.To = to
	return err
}

func (s *JSONV3SwapExactInParams) ConvertToBigIntType() (*V3SwapExactInParams, error) {
	amountIn, ok := new(big.Int).SetString(s.AmountIn, 10)
	amountOutMin, ok2 := new(big.Int).SetString(s.AmountOutMin, 10)
	if !ok || !ok2 {
		log.Warn().Msg("V3SwapExactInParams: failed to convert string to big.Int")
		return nil, fmt.Errorf("failed to convert string to big.Int")
	}
	return &V3SwapExactInParams{
		AmountIn:     amountIn,
		AmountOutMin: amountOutMin,
		Path:         s.Path,
		To:           s.To,
		PayerIsUser:  s.PayerIsUser,
	}, nil
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
