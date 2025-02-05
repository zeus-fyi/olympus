package artemis_trading_types

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
)

type TradeOutcome struct {
	AmountIn            *big.Int         `json:"amountIn"`
	AmountInAddr        accounts.Address `json:"amountInAddr"`
	AmountFees          *big.Int         `json:"amountFees"`
	AmountOut           *big.Int         `json:"amountOut"`
	AmountOutDrift      *big.Int         `json:"amountOutDrift,omitempty"`
	AmountOutAddr       accounts.Address `json:"amountOutAddr"`
	StartReservesToken0 *big.Int         `json:"startReservesToken0"`
	StartReservesToken1 *big.Int         `json:"startReservesToken1"`
	EndReservesToken0   *big.Int         `json:"endReservesToken0"`
	EndReservesToken1   *big.Int         `json:"endReservesToken1"`

	SimulatedAmountOut    *big.Int `json:"simulatedAmountOut,omitempty"`
	PreTradeTokenBalance  *big.Int `json:"preTradeTokenBalance,omitempty"`
	PostTradeTokenBalance *big.Int `json:"postTradeTokenBalance,omitempty"`
	DiffTradeTokenBalance *big.Int `json:"diffTradeTokenBalance,omitempty"`

	PostTradeEthBalance *big.Int             `json:"postTradeEthBalance,omitempty"`
	PreTradeEthBalance  *big.Int             `json:"preTradeEthBalance,omitempty"`
	DiffTradeEthBalance *big.Int             `json:"diffTradeEthBalance,omitempty"`
	OrderedTxs          []accounts.Hash      `json:"orderedTxs,omitempty"`
	BundleTxs           []*types.Transaction `json:"bundleTxs,omitempty"`
	TotalGasCost        uint64               `json:"totalGasCost,omitempty"`
}

func (t *TradeOutcome) AreTradeParamsValid() bool {
	if artemis_eth_units.AreAnyValuesLessThanEqZeroOrOne(t.AmountIn, t.AmountOut) {
		return false
	}
	return true
}

type TxLifecycleStats struct {
	TxHash     accounts.Hash
	GasUsed    uint64
	RxBlockNum uint64
}

func (t *TradeOutcome) GetGasUsageForAllTxs(ctx context.Context, w web3_actions.Web3Actions) error {
	for _, tx := range t.OrderedTxs {
		txInfo, err := GetTxLifecycleStats(ctx, w, accounts.HexToHash(tx.Hex()))
		if err != nil {
			return err
		}
		t.TotalGasCost += txInfo.GasUsed
	}
	return nil
}

func GetTxLifecycleStats(ctx context.Context, w web3_actions.Web3Actions, txHash accounts.Hash) (TxLifecycleStats, error) {
	tx, _, err := w.C.TransactionByHash(ctx, common.Hash(txHash))
	if err != nil {
		log.Err(err).Msg("GetTxLifecycleStats: error getting tx by hash")
		return TxLifecycleStats{}, err
	}
	rx, err := w.C.TransactionReceipt(ctx, common.Hash(txHash))
	if err != nil {
		log.Err(err).Msg("GetTxLifecycleStats: error getting rx by hash")
		return TxLifecycleStats{}, err
	}
	return TxLifecycleStats{
		TxHash:     txHash,
		GasUsed:    rx.GasUsed * tx.GasPrice().Uint64(),
		RxBlockNum: rx.BlockNumber.Uint64(),
	}, err
}

func (t *TradeOutcome) PrintDebug() {
	fmt.Println("amountInAddr", t.AmountInAddr.String(), "amountIn", t.AmountIn.String())
	fmt.Println("amountOutAddr", t.AmountOutAddr.String(), "amountOut", t.AmountOut.String())
}

func (t *TradeOutcome) PostTradeGasAdjustedBalance() *big.Int {
	tgCost := new(big.Int).SetUint64(t.TotalGasCost)
	return new(big.Int).Sub(t.SimulatedAmountOut, tgCost)
}

func (t *TradeOutcome) AddTxHash(tx accounts.Hash) {
	if t.OrderedTxs == nil {
		t.OrderedTxs = []accounts.Hash{}
	}
	t.OrderedTxs = append(t.OrderedTxs, tx)
}

func (t *TradeOutcome) AddTx(tx *types.Transaction) {
	if t.BundleTxs == nil {
		t.BundleTxs = []*types.Transaction{}
	}
	t.BundleTxs = append(t.BundleTxs, tx)
}

type JSONTradeOutcome struct {
	AmountIn            string           `json:"amountIn"`
	AmountInAddr        accounts.Address `json:"amountInAddr"`
	AmountFees          string           `json:"amountFees"`
	AmountOut           string           `json:"amountOut"`
	AmountOutAddr       accounts.Address `json:"amountOutAddr"`
	StartReservesToken0 string           `json:"startReservesToken0"`
	StartReservesToken1 string           `json:"startReservesToken1"`
	EndReservesToken0   string           `json:"endReservesToken0"`
	EndReservesToken1   string           `json:"endReservesToken1"`
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
