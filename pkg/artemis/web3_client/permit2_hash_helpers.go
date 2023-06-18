package web3_client

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeus-fyi/gochain/web3/accounts"
)

var (
	PermitDetailsTypeHash           = crypto.Keccak256Hash([]byte("PermitDetails(address token,uint160 amount,uint48 expiration,uint48 nonce)"))
	PermitSingleTypeHash            = crypto.Keccak256Hash([]byte("PermitSingle(PermitDetails details,address spender,uint256 sigDeadline)PermitDetails(address token,uint160 amount,uint48 expiration,uint48 nonce)"))
	TokenPermissionsTypeHash        = crypto.Keccak256Hash([]byte("TokenPermissions(address token,uint256 amount)"))
	PermitTransferFromTypeHash      = crypto.Keccak256Hash([]byte("PermitTransferFrom(TokenPermissions permitted,address spender,uint256 nonce,uint256 deadline)TokenPermissions(address token,uint256 amount)"))
	PermitBatchTypeHash             = crypto.Keccak256Hash([]byte("PermitBatch(PermitDetails[] details,address spender,uint256 sigDeadline)PermitDetails(address token,uint160 amount,uint48 expiration,uint48 nonce)"))
	PermitBatchTransferFromTypeHash = crypto.Keccak256Hash([]byte("PermitBatchTransferFrom(TokenPermissions[] permitted,address spender,uint256 nonce,uint256 deadline)TokenPermissions(address token,uint256 amount)"))
	TokenPermissionsTypeString      = "TokenPermissions(address token,uint256 amount)"
)

func _hashPermitDetails(permitDetails PermitDetails) []byte {
	parsedABI, err := abi.JSON(strings.NewReader(`[{"type":"function","inputs":[{"type":"bytes32"},{"type":"tuple","components":[{"name":"token","type":"address"},{"name":"amount","type":"uint160"},{"name":"expiration","type":"uint48"},{"name":"nonce","type":"uint48"}],"name":"details"}],"name":"abiEncode"}]`))
	if err != nil {
		panic(err)
	}
	data, err := parsedABI.Methods["abiEncode"].Inputs.Pack(common.BytesToHash(PermitDetailsTypeHash.Bytes()), permitDetails)
	if err != nil {
		panic(err)
	}
	hashed := crypto.Keccak256Hash(data)
	return hashed.Bytes()
}

func hashPermitSingle(permitSingle PermitSingle) common.Hash {
	hashed := _hashPermitDetails(permitSingle.PermitDetails)
	parsedABI, err := abi.JSON(strings.NewReader(`[{"type":"function","inputs":[{"type":"bytes32"},{"type":"bytes32"},{"type":"address"},{"type":"uint256"}],"name":"abiEncode","outputs":[{"type":"bytes"}]}]`))
	if err != nil {
		panic(err)
	}

	data, err := parsedABI.Methods["abiEncode"].Inputs.Pack(common.BytesToHash(PermitSingleTypeHash.Bytes()), common.BytesToHash(hashed), permitSingle.Spender, permitSingle.SigDeadline)
	if err != nil {
		panic(err)
	}
	return crypto.Keccak256Hash(data)
}

func _hashTokenPermissions(permitted TokenPermissions) []byte {
	parsedABI, err := abi.JSON(strings.NewReader(`[{"type":"function","inputs":[{"type":"bytes32"},{"type":"address"},{"type":"uint256"}],"name":"abiEncode","outputs":[{"type":"bytes"}]}]`))
	if err != nil {
		panic(err)
	}
	data, err := parsedABI.Methods["abiEncode"].Inputs.Pack(TokenPermissionsTypeHash, permitted.Token, permitted.Amount)
	if err != nil {
		panic(err)
	}
	hashed := crypto.Keccak256(data)
	return hashed
}

func hashPermitTransferFrom(permitTransferFrom PermitTransferFrom, sender accounts.Address) common.Hash {
	tokenPermissions := _hashTokenPermissions(permitTransferFrom.TokenPermissions)
	parsedABI, err := abi.JSON(strings.NewReader(`[{"type":"function","inputs":[{"type":"bytes32"},{"type":"bytes32"},{"type":"address"},{"type":"uint256"},{"type":"uint256"}],"name":"abiEncode"}]`))
	if err != nil {
		panic(err)
	}
	data, err := parsedABI.Methods["abiEncode"].Inputs.Pack(common.BytesToHash(PermitTransferFromTypeHash.Bytes()), common.BytesToHash(tokenPermissions), sender, permitTransferFrom.Nonce, permitTransferFrom.SigDeadline)
	if err != nil {
		panic(err)
	}
	return crypto.Keccak256Hash(data)
}
