package web3_client

import (
	"context"
	"math/big"

	"github.com/zeus-fyi/gochain/v4/common"
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

	SimulatedAmountOut  *big.Int      `json:"simulatedAmountOut,omitempty"`
	PostTradeEthBalance *big.Int      `json:"postTradeEthBalance,omitempty"`
	PreTradeEthBalance  *big.Int      `json:"preTradeEthBalance,omitempty"`
	OrderedTxs          []common.Hash `json:"orderedTxs,omitempty"`
	TotalGasCost        uint64        `json:"totalGasCost,omitempty"`
}

func (t *TradeOutcome) PostTradeGasAdjustedBalance() *big.Int {
	tgCost := new(big.Int).SetUint64(t.TotalGasCost)
	return new(big.Int).Sub(t.SimulatedAmountOut, tgCost)
}

func (t *TradeOutcome) AddTxHash(tx common.Hash) {
	if t.OrderedTxs == nil {
		t.OrderedTxs = []common.Hash{}
	}
	t.OrderedTxs = append(t.OrderedTxs, tx)
}

func (t *TradeOutcome) GetGasUsageForAllTxs(ctx context.Context, w Web3Client) error {
	for _, tx := range t.OrderedTxs {
		txInfo, err := w.GetTxLifecycleStats(ctx, tx)
		if err != nil {
			return err
		}
		t.TotalGasCost += txInfo.GasUsed
	}
	return nil
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
