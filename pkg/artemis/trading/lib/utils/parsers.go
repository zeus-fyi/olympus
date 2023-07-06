package artemis_utils

import (
	"fmt"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
)

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

func StringsToAddresses(addressOne, addressTwo string) (accounts.Address, accounts.Address) {
	addrOne := accounts.HexToAddress(addressOne)
	addrTwo := accounts.HexToAddress(addressTwo)
	return addrOne, addrTwo
}
