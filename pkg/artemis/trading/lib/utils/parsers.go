package artemis_utils

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"
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

func StringsToAddresses(addressOne, addressTwo string) (accounts.Address, accounts.Address) {
	addrOne := accounts.HexToAddress(addressOne)
	addrTwo := accounts.HexToAddress(addressTwo)
	return addrOne, addrTwo
}
