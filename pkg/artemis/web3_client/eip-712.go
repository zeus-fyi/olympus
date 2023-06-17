package web3_client

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

type EIP712 struct {
	contractAddress       accounts.Address
	cachedDomainSeparator common.Hash
	cachedChainID         *big.Int
	hashedName            common.Hash
	typeHash              common.Hash
}

func NewEIP712(chainID *big.Int, contractAddress accounts.Address, name string) *EIP712 {
	hashedName := crypto.Keccak256Hash([]byte(name))
	typeHash := crypto.Keccak256Hash([]byte("EIP712Domain(string name,uint256 chainId,address verifyingContract)"))

	return &EIP712{
		cachedChainID:         chainID,
		cachedDomainSeparator: buildDomainSeparator(typeHash, hashedName, chainID, contractAddress),
		hashedName:            hashedName,
		typeHash:              typeHash,
	}
}

func NewEIP712ForPermit2(chainID *big.Int, contractAddress accounts.Address) *EIP712 {
	return NewEIP712(chainID, contractAddress, "Permit2")
}

func buildDomainSeparator(typeHash common.Hash, nameHash common.Hash, chainID *big.Int, contractAddress accounts.Address) common.Hash {
	parsedABI, err := abi.JSON(strings.NewReader(`[{"type":"function","inputs":[{"type":"bytes32"},{"type":"bytes32"},{"type":"uint256"},{"type":"address"}],"name":"abiEncode"}]`))
	if err != nil {
		log.Err(err)
		panic(err)
	}
	data, err := parsedABI.Methods["abiEncode"].Inputs.Pack(typeHash, nameHash, chainID, contractAddress)
	if err != nil {
		log.Err(err)
		panic(err)
	}
	return crypto.Keccak256Hash(data)
}

func (e *EIP712) DomainSeparator() common.Hash {
	return e.cachedDomainSeparator
}

func (e *EIP712) HashTypedData(dataHash common.Hash) common.Hash {
	domainSeparator := e.DomainSeparator()
	encodedData := append([]byte("\x19\x01"), domainSeparator.Bytes()...)
	encodedData = append(encodedData, dataHash.Bytes()...)
	return crypto.Keccak256Hash(encodedData)
}
