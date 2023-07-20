package web3_client

import (
	"errors"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

type SandwichTradePrediction struct {
	SellAmount     *big.Int `json:"sellAmount"`
	ExpectedProfit *big.Int `json:"expectedProfit"`
}

type JSONSandwichTradePrediction struct {
	SellAmount     string `json:"sellAmount"`
	ExpectedProfit string `json:"expectedProfit"`
}

func (s *JSONSandwichTradePrediction) ConvertToBigIntType() (SandwichTradePrediction, error) {
	sellAmount, ok1 := new(big.Int).SetString(s.SellAmount, 10)
	expectedProfit, ok2 := new(big.Int).SetString(s.ExpectedProfit, 10)
	if !ok1 || !ok2 {
		log.Err(errors.New("error converting to big int")).Msg("JSONSandwichTradePrediction: error converting to big int")
		return SandwichTradePrediction{}, errors.New("error converting to big int")
	}
	return SandwichTradePrediction{
		SellAmount:     sellAmount,
		ExpectedProfit: expectedProfit,
	}, nil
}

func (s *SandwichTradePrediction) CheckForValidityAndProfit() bool {
	return s.ExpectedProfit.Cmp(big.NewInt(0)) > 0 && s.SellAmount.Cmp(big.NewInt(0)) > 0
}

func (s *SandwichTradePrediction) ConvertToJSONType() JSONSandwichTradePrediction {
	if s.SellAmount == nil {
		s.SellAmount = big.NewInt(0)
	}
	if s.ExpectedProfit == nil {
		s.ExpectedProfit = big.NewInt(0)
	}
	return JSONSandwichTradePrediction{
		SellAmount:     s.SellAmount.String(),
		ExpectedProfit: s.ExpectedProfit.String(),
	}
}

type AddLiquidityParams struct {
	TokenA         accounts.Address `json:"tokenA"`
	TokenB         accounts.Address `json:"tokenB"`
	AmountADesired *big.Int         `json:"amountADesired"`
	AmountBDesired *big.Int         `json:"amountBDesired"`
	AmountAMin     *big.Int         `json:"amountAMin"`
	AmountBMin     *big.Int         `json:"amountBMin"`
	To             accounts.Address `json:"to"`
	Deadline       *big.Int         `json:"deadline"`
}

type AddLiquidityETHParams struct {
	Token              accounts.Address `json:"token"`
	AmountTokenDesired *big.Int         `json:"amountTokenDesired"`
	AmountTokenMin     *big.Int         `json:"amountTokenMin"`
	AmountETHMin       *big.Int         `json:"amountETHMin"`
	To                 accounts.Address `json:"to"`
	Deadline           *big.Int         `json:"deadline"`
}

type RemoveLiquidityParams struct {
	TokenA     accounts.Address `json:"tokenA"`
	TokenB     accounts.Address `json:"tokenB"`
	Liquidity  *big.Int         `json:"liquidity"`
	AmountAMin *big.Int         `json:"amountAMin"`
	AmountBMin *big.Int         `json:"amountBMin"`
	To         accounts.Address `json:"to"`
	Deadline   *big.Int         `json:"deadline"`
}

type RemoveLiquidityETHParams struct {
	Token          accounts.Address `json:"token"`
	Liquidity      *big.Int         `json:"liquidity"`
	AmountTokenMin *big.Int         `json:"amountTokenMin"`
	AmountETHMin   *big.Int         `json:"amountETHMin"`
	To             accounts.Address `json:"to"`
	Deadline       *big.Int         `json:"deadline"`
}

type RemoveLiquidityWithPermitParams struct {
	TokenA     accounts.Address `json:"tokenA"`
	TokenB     accounts.Address `json:"tokenB"`
	Liquidity  *big.Int         `json:"liquidity"`
	AmountAMin *big.Int         `json:"amountAMin"`
	AmountBMin *big.Int         `json:"amountBMin"`
	To         accounts.Address `json:"to"`
	Deadline   *big.Int         `json:"deadline"`
	ApproveMax bool             `json:"approveMax"`
	V          uint8            `json:"v"`
	R          [32]byte         `json:"r"`
	S          [32]byte         `json:"s"`
}

type RemoveLiquidityETHWithPermitParams struct {
	Token          accounts.Address `json:"token"`
	Liquidity      *big.Int         `json:"liquidity"`
	AmountTokenMin *big.Int         `json:"amountTokenMin"`
	AmountETHMin   *big.Int         `json:"amountETHMin"`
	To             accounts.Address `json:"to"`
	Deadline       *big.Int         `json:"deadline"`
	ApproveMax     bool             `json:"approveMax"`
	V              uint8            `json:"v"`
	R              [32]byte         `json:"r"`
	S              [32]byte         `json:"s"`
}
