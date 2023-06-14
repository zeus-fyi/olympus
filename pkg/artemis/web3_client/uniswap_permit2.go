package web3_client

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/zeus-fyi/gochain/web3/accounts"
)

/*
PERMIT2_PERMIT
IAllowanceTransfer.PermitSingle A PermitSingle struct outlining the Permit2 permit to execute
bytes The signature to provide to Permit2

The individual that signed the permit must be the msg.sender of the transaction

PERMIT2_TRANSFER_FROM
address The token to fetch from Permit2
address The recipient of the tokens fetched
uint256 The amount of token to fetch
The individual that the tokens are fetched from is always the msg.sender of the transaction

PERMIT2_PERMIT_BATCH
IAllowanceTransfer.PermitBatch A PermitBatch struct outlining all of the Permit2 permits to execute.
bytes The signature to provide to Permit2
The individual that signed the permits must be the msg.sender of the transaction
*/

const (
	Permit2TransferFrom      = "PERMIT2_TRANSFER_FROM"
	Permit2PermitBatch       = "PERMIT2_PERMIT_BATCH"
	Permit2Permit            = "PERMIT2_PERMIT"
	Permit2TransferFromBatch = "PERMIT2_TRANSFER_FROM_BATCH"

	Permit2SmartContractAddress = "0x000000000022D473030F116dDEE9F6B43aC78BA3"
)

var Permit2AbiDecoder = MustLoadPermit2Abi()

type Permit2TransferFromParams struct {
	Token     accounts.Address `json:"token"`
	Recipient accounts.Address `json:"recipient"`
	Amount    *big.Int         `json:"amount"`
}

func (p *Permit2TransferFromParams) Encode(ctx context.Context) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[Permit2TransferFrom].Inputs.Pack(p.Token, p.Recipient, p.Amount)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (p *Permit2TransferFromParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[Permit2TransferFrom].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	token, err := ConvertToAddress(args["token"])
	if err != nil {
		return err
	}
	recipient, err := ConvertToAddress(args["recipient"])
	if err != nil {
		return err
	}
	amount, err := ParseBigInt(args["amount"])
	if err != nil {
		return err
	}
	p.Token = token
	p.Recipient = recipient
	p.Amount = amount
	return nil
}

type Permit2PermitParams struct {
	PermitSingle
	Signature []byte `json:"signature"`
}

func (p *Permit2PermitParams) Sign(acc *accounts.Account, chainID *big.Int, contractAddress accounts.Address) error {
	if acc == nil {
		return errors.New("account is nil")
	}
	hashed := hashPermitSingle(p.PermitSingle)
	eip := NewEIP712ForPermit2(chainID, contractAddress)
	hashed = eip.HashTypedData(hashed)
	sig, err := acc.Sign(hashed.Bytes())
	if err != nil {
		return err
	}
	p.Signature = sig
	return nil
}

// equivalent: abi.decode(inputs, (IAllowanceTransfer.PermitSingle, bytes))

type PermitTransferFrom struct {
	TokenPermissions
	Expiration  *big.Int `json:"expiration"`  // uint48 can be represented as uint64 in Go
	Nonce       *big.Int `json:"nonce"`       // uint48 can be represented as uint64 in Go
	SigDeadline *big.Int `json:"sigDeadline"` // uint48 can be represented as uint64 in Go
}

type PermitSingle struct {
	PermitDetails
	Spender     accounts.Address `json:"spender"`
	SigDeadline *big.Int         `json:"sigDeadline"` // uint48 can be represented as uint64 in Go
}

type TokenPermissions struct {
	Token  accounts.Address `json:"token"`
	Amount *big.Int         `json:"amount"` // uint160 can be represented as *big.Int in Go
}

type PermitDetails struct {
	TokenPermissions
	Expiration *big.Int `json:"expiration"` // uint48 can be represented as uint64 in Go
	Nonce      *big.Int `json:"nonce"`      // uint48 can be represented as uint64 in Go
}

func (p *Permit2PermitParams) Encode(ctx context.Context) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[Permit2Permit].Inputs.Pack(p.Token, p.Amount, p.Expiration, p.Nonce, p.Spender, p.SigDeadline, p.Signature)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}
func (p *Permit2PermitParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[Permit2Permit].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	token, err := ConvertToAddress(args["token"])
	if err != nil {
		return err
	}
	amount, err := ParseBigInt(args["amount"])
	if err != nil {
		return err
	}
	expiration, err := ParseBigInt(args["expiration"])
	if err != nil {
		return err
	}
	nonce, err := ParseBigInt(args["nonce"])
	if err != nil {
		return err
	}
	spender, err := ConvertToAddress(args["spender"])
	if err != nil {
		return err
	}
	sigDeadline, err := ParseBigInt(args["sigDeadline"])
	if err != nil {
		return err
	}
	signature := args["signature"].([]byte)
	p.Token = token
	p.Amount = amount
	p.Expiration = expiration
	p.Nonce = nonce
	p.Spender = spender
	p.SigDeadline = sigDeadline
	p.Signature = signature
	return nil
}

type Permit2PermitBatchParams struct {
	PermitBatch PermitBatch `json:"permitBatch"`
	Signature   []byte      `json:"signature"`
}

type PermitBatch struct {
	Details     []PermitDetails  `json:"details"`
	Spender     accounts.Address `json:"spender"`
	SigDeadline *big.Int         `json:"sigDeadline"`
}

// abi.decode(inputs, (IAllowanceTransfer.PermitBatch, bytes));

func (p *Permit2PermitBatchParams) Encode(ctx context.Context) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[Permit2PermitBatch].Inputs.Pack(p.PermitBatch, p.Signature)
	if err != nil {
		return nil, err
	}

	return inputs, nil
}

func (p *Permit2PermitBatchParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[Permit2PermitBatch].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	b, err := json.Marshal(args["permitBatch"])
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &p.PermitBatch)
	if err != nil {
		return err
	}
	signature := args["signature"].([]byte)
	p.Signature = signature
	return nil
}

type Permit2PermitTransferFromBatchParams struct {
	Details []AllowanceTransferDetails `json:"batchDetails"`
}

type AllowanceTransferDetails struct {
	From   accounts.Address `json:"from"`
	To     accounts.Address `json:"to"`
	Amount *big.Int         `json:"amount"`
	Token  accounts.Address `json:"token"`
}

// abi.decode(inputs, (IAllowanceTransfer.AllowanceTransferDetails[]));

func (p *Permit2PermitTransferFromBatchParams) Encode(ctx context.Context) ([]byte, error) {
	inputs, err := UniversalRouterDecoderAbi.Methods[Permit2TransferFromBatch].Inputs.Pack(p.Details)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (p *Permit2PermitTransferFromBatchParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoderAbi.Methods[Permit2TransferFromBatch].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	b, err := json.Marshal(args["batchDetails"])
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &p.Details)
	if err != nil {
		return err
	}
	return nil
}
