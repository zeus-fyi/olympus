package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeus-fyi/gochain/web3/accounts"
	entities "github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_core/entities"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_libs/uniswap_v3/constants"
)

/**
 * Computes a pool address
 * @param factoryAddress The Uniswap V3 factory address
 * @param tokenA The first token of the pair, irrespective of sort order
 * @param tokenB The second token of the pair, irrespective of sort order
 * @param fee The fee tier of the pool
 * @returns The pool address
 */

func ComputePoolAddress(factoryAddress accounts.Address, tokenA *entities.Token, tokenB *entities.Token, fee constants.FeeAmount, initCodeHashManualOverride string) (accounts.Address, error) {
	isSorted, err := tokenA.SortsBefore(tokenB)
	if err != nil {
		return accounts.Address{}, err
	}
	var (
		token0 *entities.Token
		token1 *entities.Token
	)
	if isSorted {
		token0 = tokenA
		token1 = tokenB
	} else {
		token0 = tokenB
		token1 = tokenA
	}
	return getCreate2Address(factoryAddress, token0.Address, token1.Address, fee, initCodeHashManualOverride), nil
}

func getCreate2Address(factoyAddress, addressA, addressB accounts.Address, fee constants.FeeAmount, initCodeHashManualOverride string) accounts.Address {
	var salt [32]byte
	copy(salt[:], crypto.Keccak256(abiEncode(addressA, addressB, fee)))

	if initCodeHashManualOverride != "" {
		crypto.CreateAddress2(common.HexToAddress(factoyAddress.Hex()), salt, common.FromHex(initCodeHashManualOverride))
	}
	return accounts.HexToAddress(crypto.CreateAddress2(common.HexToAddress(factoyAddress.Hex()), salt, common.FromHex(constants.PoolInitCodeHash)).Hex())
}

func abiEncode(addressA, addressB accounts.Address, fee constants.FeeAmount) []byte {
	addressTy, _ := abi.NewType("address", "address", nil)
	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)

	arguments := abi.Arguments{{Type: addressTy}, {Type: addressTy}, {Type: uint256Ty}}

	bytes, _ := arguments.Pack(
		addressA,
		addressB,
		big.NewInt(int64(fee)),
	)
	return bytes
}
