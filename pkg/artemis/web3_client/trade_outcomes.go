package web3_client

import (
	"math/big"

	"github.com/gochain/gochain/v4/common"
)

type TradeOutcome struct {
	AmountIn            *big.Int       `json:"amountIn"`
	AmountInAddr        common.Address `json:"amountInAddr"`
	AmountFees          *big.Int       `json:"amountFees"`
	AmountOut           *big.Int       `json:"amountOut"`
	AmountOutAddr       common.Address `json:"amountOutAddr"`
	StartReservesToken0 *big.Int       `json:"startReservesToken0"`
	StartReservesToken1 *big.Int       `json:"startReservesToken1"`
	EndReservesToken0   *big.Int       `json:"endReservesToken0"`
	EndReservesToken1   *big.Int       `json:"endReservesToken1"`

	OrderedTxs   []common.Hash `json:"orderedTxs,omitempty"`
	TotalGasCost *big.Int      `json:"totalGasCost,omitempty"`
}

func (t *TradeOutcome) AddTxHash(tx common.Hash) {
	if t.OrderedTxs == nil {
		t.OrderedTxs = []common.Hash{}
	}
	t.OrderedTxs = append(t.OrderedTxs, tx)
}

type JSONTradeOutcome struct {
	AmountIn            string         `json:"amountIn"`
	AmountInAddr        common.Address `json:"amountInAddr"`
	AmountFees          string         `json:"amountFees"`
	AmountOut           string         `json:"amountOut"`
	AmountOutAddr       common.Address `json:"amountOutAddr"`
	StartReservesToken0 string         `json:"startReservesToken0"`
	StartReservesToken1 string         `json:"startReservesToken1"`
	EndReservesToken0   string         `json:"endReservesToken0"`
	EndReservesToken1   string         `json:"endReservesToken1"`
}

func (t *JSONTradeOutcome) ConvertToBigIntType() TradeOutcome {
	amountIn, _ := new(big.Int).SetString(t.AmountIn, 10)
	amountFees, _ := new(big.Int).SetString(t.AmountFees, 10)
	amountOut, _ := new(big.Int).SetString(t.AmountOut, 10)
	startReservesToken0, _ := new(big.Int).SetString(t.StartReservesToken0, 10)
	startReservesToken1, _ := new(big.Int).SetString(t.StartReservesToken1, 10)
	endReservesToken0, _ := new(big.Int).SetString(t.EndReservesToken0, 10)
	endReservesToken1, _ := new(big.Int).SetString(t.EndReservesToken1, 10)
	return TradeOutcome{
		AmountIn:            amountIn,
		AmountInAddr:        t.AmountInAddr,
		AmountFees:          amountFees,
		AmountOut:           amountOut,
		AmountOutAddr:       t.AmountOutAddr,
		StartReservesToken0: startReservesToken0,
		StartReservesToken1: startReservesToken1,
		EndReservesToken0:   endReservesToken0,
		EndReservesToken1:   endReservesToken1,
	}
}

func (t *TradeOutcome) ConvertToJSONType() JSONTradeOutcome {
	return JSONTradeOutcome{
		AmountIn:            t.AmountIn.String(),
		AmountInAddr:        t.AmountInAddr,
		AmountFees:          t.AmountFees.String(),
		AmountOut:           t.AmountOut.String(),
		AmountOutAddr:       t.AmountOutAddr,
		StartReservesToken0: t.StartReservesToken0.String(),
		StartReservesToken1: t.StartReservesToken1.String(),
		EndReservesToken0:   t.EndReservesToken0.String(),
		EndReservesToken1:   t.EndReservesToken1.String(),
	}
}
