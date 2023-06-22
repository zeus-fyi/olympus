package web3_client

import (
	"context"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
)

const (
	V3SwapExactIn = "V3_SWAP_EXACT_IN"
)

type V3SwapExactInParams struct {
	AmountIn     *big.Int         `json:"amountIn"`
	AmountOutMin *big.Int         `json:"amountOutMin"`
	Path         TokenFeePath     `json:"path"`
	To           accounts.Address `json:"to"`
	PayerIsUser  bool             `json:"payerIsUser"`
}

type JSONV3SwapExactInParams struct {
	AmountIn     string           `json:"amountIn"`
	AmountOutMin string           `json:"amountOutMin"`
	Path         TokenFeePath     `json:"path"`
	To           accounts.Address `json:"to"`
	PayerIsUser  bool             `json:"payerIsUser"`
}

func (s *V3SwapExactInParams) Encode(ctx context.Context) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[V3SwapExactIn].Inputs.Pack(s.To, s.AmountIn, s.AmountOutMin, s.Path.Encode(), s.PayerIsUser)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (s *V3SwapExactInParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[V3SwapExactIn].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	amountIn, err := ParseBigInt(args["amountIn"])
	if err != nil {
		return err
	}
	amountOutMin, err := ParseBigInt(args["amountOutMin"])
	if err != nil {
		return err
	}
	pathBytes := args["path"].([]byte)
	hexStr := accounts.Bytes2Hex(pathBytes)
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
	tfp.Reverse()

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
