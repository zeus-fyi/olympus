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
	V3SwapExactOut = "V3_SWAP_EXACT_OUT"
)

type V3SwapExactOutParams struct {
	AmountInMax *big.Int                           `json:"amountInMax"`
	AmountOut   *big.Int                           `json:"amountOut"`
	Path        artemis_trading_types.TokenFeePath `json:"path"`
	To          accounts.Address                   `json:"to"`
	PayerIsUser bool                               `json:"payerIsUser"`
}

type JSONV3SwapExactOutParams struct {
	AmountInMax string                             `json:"amountInMax"`
	AmountOut   string                             `json:"amountOut"`
	Path        artemis_trading_types.TokenFeePath `json:"path"`
	To          accounts.Address                   `json:"to"`
	PayerIsUser bool                               `json:"payerIsUser"`
}

func (s *V3SwapExactOutParams) Encode(ctx context.Context, abiFile *abi.ABI) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[V3SwapExactOut].Inputs.Pack(s.To, s.AmountOut, s.AmountInMax, s.Path.Encode(), s.PayerIsUser)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (s *V3SwapExactOutParams) Decode(ctx context.Context, data []byte, abiFile *abi.ABI) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[V3SwapExactOut].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		log.Warn().Err(err).Msg("V3SwapExactOutParams: UniversalRouterDecoderAbi failed to unpack data")
		return err
	}

	amountInMax, err := ParseBigInt(args["amountInMax"])
	if err != nil {
		log.Warn().Err(err).Msg("V3SwapExactOutParams: failed to parse amountInMax")
		return err
	}
	amountOut, err := ParseBigInt(args["amountOut"])
	if err != nil {
		log.Warn().Err(err).Msg("V3SwapExactOutParams: failed to parse amountOut")
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
			log.Warn().Err(err).Msg("V3SwapExactOutParams: failed to parse fee")
			return err
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
		log.Warn().Err(err).Msg("V3SwapExactOutParams: failed to parse recipient")
		return err
	}
	payerIsUserInterface, pok := args["payerIsUser"]
	if !pok || payerIsUserInterface == nil {
		// Handle the situation when args["path"] doesn't exist or is nil
		log.Warn().Msg("V3SwapExactOutParams: 'payerIsUser' does not exist or is nil")
		return fmt.Errorf("payerIsUser does not exist or is nil")
	}
	payerIsUserBool, ok := payerIsUserInterface.(bool)
	if !ok {
		// Handle the situation when the conversion fails
		log.Warn().Msg("V3SwapExactOutParams: failed to convert 'payerIsUser' to bool")
		return fmt.Errorf("failed to convert payerIsUser to bool")
	}
	s.AmountInMax = amountInMax
	s.AmountOut = amountOut
	s.Path = tfp
	s.To = to
	s.PayerIsUser = payerIsUserBool
	return nil
}

func (s *JSONV3SwapExactOutParams) ConvertToBigIntType() (*V3SwapExactOutParams, error) {
	amountInMax, ok := new(big.Int).SetString(s.AmountInMax, 10)
	amountOut, ok1 := new(big.Int).SetString(s.AmountOut, 10)
	if !ok || !ok1 {
		log.Warn().Msg("V3SwapExactOutParams: failed to convert string to big.Int")
		return nil, errors.New("failed to convert string to big.Int")
	}
	return &V3SwapExactOutParams{
		AmountInMax: amountInMax,
		AmountOut:   amountOut,
		Path:        s.Path,
		To:          s.To,
		PayerIsUser: s.PayerIsUser,
	}, nil
}

func (s *V3SwapExactOutParams) ConvertToJSONType() *JSONV3SwapExactOutParams {
	return &JSONV3SwapExactOutParams{
		AmountInMax: s.AmountInMax.String(),
		AmountOut:   s.AmountOut.String(),
		Path:        s.Path,
		To:          s.To,
		PayerIsUser: s.PayerIsUser,
	}
}
