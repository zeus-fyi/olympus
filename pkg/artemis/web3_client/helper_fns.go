package web3_client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/gochain/gochain/v4/common"
	"github.com/rs/zerolog/log"
)

func StringsToAddresses(addressOne, addressTwo string) (common.Address, common.Address) {
	addrOne := common.HexToAddress(addressOne)
	addrTwo := common.HexToAddress(addressTwo)
	return addrOne, addrTwo
}

type TxLifecycleStats struct {
	TxHash     common.Hash
	GasUsed    uint64
	TxBlockNum *big.Int
	RxBlockNum uint64
}

func (w *Web3Client) GetTxLifecycleStats(ctx context.Context, txHash common.Hash) (TxLifecycleStats, error) {
	tx, err := w.GetTransactionByHash(ctx, txHash)
	if err != nil {
		log.Err(err).Msg("GetTxLifecycleStats: error getting tx by hash")
		return TxLifecycleStats{}, err
	}
	rx, err := w.GetTransactionReceipt(ctx, txHash)
	if err != nil {
		log.Err(err).Msg("GetTxLifecycleStats: error getting rx by hash")
		return TxLifecycleStats{}, err
	}
	return TxLifecycleStats{
		TxHash:     txHash,
		GasUsed:    rx.GasUsed * tx.GasPrice.Uint64(),
		TxBlockNum: tx.BlockNumber,
		RxBlockNum: rx.BlockNumber,
	}, err
}

func (w *Web3Client) GetEthBalance(ctx context.Context, addr string, blockNum *big.Int) (*big.Int, error) {
	w.Dial()
	defer w.Close()
	balance, err := w.GetBalance(ctx, addr, blockNum)
	if err != nil {
		return balance, err
	}
	return balance, err
}

func ConvertAmountsToBigIntSlice(amounts []interface{}) []*big.Int {
	var amountsBigInt []*big.Int
	for _, amount := range amounts {
		pair := amount.([]*big.Int)
		for _, p := range pair {
			amountsBigInt = append(amountsBigInt, p)
		}
	}
	return amountsBigInt
}

func ParseBigInt(i interface{}) (*big.Int, error) {
	switch v := i.(type) {
	case *big.Int:
		return i.(*big.Int), nil
	case string:
		base := 10
		result := new(big.Int)
		_, ok := result.SetString(v, base)
		if !ok {
			return nil, fmt.Errorf("failed to parse string '%s' into big.Int", v)
		}
		return result, nil
	case uint32:
		return big.NewInt(int64(v)), nil
	case int64:
		return big.NewInt(v), nil
	default:
		return nil, fmt.Errorf("input is not a string or int64")
	}
}

func ConvertToAddressSlice(i interface{}) ([]common.Address, error) {
	switch v := i.(type) {
	case []common.Address:
		return i.([]common.Address), nil
	default:
		fmt.Println(v)
		return nil, fmt.Errorf("input is not a []common.Address")
	}
}

func ConvertToAddress(i interface{}) (common.Address, error) {
	switch v := i.(type) {
	case common.Address:
		return i.(common.Address), nil
	default:
		fmt.Println(v)
		return common.Address{}, fmt.Errorf("input is not a  common.Address")
	}
}
