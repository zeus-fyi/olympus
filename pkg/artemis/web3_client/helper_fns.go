package web3_client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

func StringsToAddresses(addressOne, addressTwo string) (accounts.Address, accounts.Address) {
	addrOne := accounts.HexToAddress(addressOne)
	addrTwo := accounts.HexToAddress(addressTwo)
	return addrOne, addrTwo
}

type TxLifecycleStats struct {
	TxHash     accounts.Hash
	GasUsed    uint64
	RxBlockNum uint64
}

func (w *Web3Client) GetTxLifecycleStats(ctx context.Context, txHash accounts.Hash) (TxLifecycleStats, error) {
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
	case int:
		return big.NewInt(int64(v)), nil
	case uint64:
		return new(big.Int).SetUint64(v), nil
	case uint32:
		return big.NewInt(int64(v)), nil
	case int64:
		return big.NewInt(v), nil
	case []byte:
		return new(big.Int).SetBytes(v), nil
	case nil:
		return big.NewInt(0), nil
	default:
		log.Warn().Msgf("ParseBigInt: unknown type %T", v)
		return big.NewInt(0), nil
	}
}

func ConvertToAddressSlice(i interface{}) ([]accounts.Address, error) {
	switch v := i.(type) {
	case []accounts.Address:
		return i.([]accounts.Address), nil
	case []common.Address:
		m := make([]accounts.Address, len(i.([]common.Address)))
		for ind, addr := range i.([]common.Address) {
			m[ind] = accounts.HexToAddress(addr.Hex())
		}
		return m, nil
	case []string:
		m := make([]accounts.Address, len(i.([]string)))
		for ind, addr := range i.([]string) {
			m[ind] = accounts.HexToAddress(addr)
		}
	case nil:
		return []accounts.Address{}, nil
	default:
		log.Warn().Msgf("ConvertToAddressSlice: unknown type %T", v)
		return nil, fmt.Errorf("input is not a []common.Address")
	}
	return nil, fmt.Errorf("input is not a []common.Address")
}

func ConvertToAddress(i interface{}) (accounts.Address, error) {
	switch v := i.(type) {
	case common.Address:
		return accounts.Address(v), nil
	case accounts.Address:
		addr := i.(common.Address)
		return accounts.HexToAddress(addr.Hex()), nil
	case nil:
		return accounts.Address{}, nil
	default:
		log.Warn().Msgf("ConvertToAddress: unknown type %T", v)
		return accounts.Address{}, fmt.Errorf("input is not a  common.Address")
	}
}
