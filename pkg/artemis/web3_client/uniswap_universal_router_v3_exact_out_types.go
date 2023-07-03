package web3_client

import (
	"context"
	"math/big"

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

func (s *V3SwapExactOutParams) Encode(ctx context.Context) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[V3SwapExactOut].Inputs.Pack(s.To, s.AmountOut, s.AmountInMax, s.Path.Encode(), s.PayerIsUser)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (s *V3SwapExactOutParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[V3SwapExactOut].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	amountInMax, err := ParseBigInt(args["amountInMax"])
	if err != nil {
		return err
	}
	amountOut, err := ParseBigInt(args["amountOut"])
	if err != nil {
		return err
	}
	pathBytes := args["path"].([]byte)
	hexStr := accounts.Bytes2Hex(pathBytes)
	tfp := artemis_trading_types.TokenFeePath{
		TokenIn: accounts.HexToAddress(hexStr[:40]),
	}
	var pathList []artemis_trading_types.TokenFee
	for i := 0; i < len(hexStr[40:]); i += 46 {
		fee, _ := new(big.Int).SetString(hexStr[40:][i:i+6], 16)
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
	payerIsUser := args["payerIsUser"].(bool)
	s.AmountInMax = amountInMax
	s.AmountOut = amountOut
	s.Path = tfp
	s.To = to
	s.PayerIsUser = payerIsUser
	return nil
}

func (s *JSONV3SwapExactOutParams) ConvertToBigIntType() *V3SwapExactOutParams {
	amountInMax, _ := new(big.Int).SetString(s.AmountInMax, 10)
	amountOut, _ := new(big.Int).SetString(s.AmountOut, 10)
	return &V3SwapExactOutParams{
		AmountInMax: amountInMax,
		AmountOut:   amountOut,
		Path:        s.Path,
		To:          s.To,
		PayerIsUser: s.PayerIsUser,
	}
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
