package web3_client

import "github.com/gochain/gochain/v4/common"

func StringsToAddresses(addressOne, addressTwo string) (common.Address, common.Address) {
	addrOne := common.HexToAddress(addressOne)
	addrTwo := common.HexToAddress(addressTwo)
	return addrOne, addrTwo
}
