package web3_client

import (
	"context"
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
)

type Permit2PermitTransferFromParams struct {
	Token     accounts.Address `json:"token"`
	Recipient accounts.Address `json:"recipient"`
	Amount    *big.Int         `json:"amount"`
}

func (p *Permit2PermitTransferFromParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoder.Methods[Permit2TransferFrom].Inputs.UnpackIntoMap(args, data)
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

// equivalent: abi.decode(inputs, (IAllowanceTransfer.PermitSingle, bytes))

type PermitSingle struct {
	PermitDetails
	Spender     accounts.Address
	SigDeadline *big.Int
}

type PermitDetails struct {
	Token      accounts.Address
	Amount     *big.Int // uint160 can be represented as *big.Int in Go
	Expiration *big.Int // uint48 can be represented as uint64 in Go
	Nonce      *big.Int // uint48 can be represented as uint64 in Go
}

func (p *Permit2PermitParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoder.Methods[Permit2Permit].Inputs.UnpackIntoMap(args, data)
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
}

// abi.decode(inputs, (IAllowanceTransfer.PermitBatch, bytes));

func (p *Permit2PermitBatchParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoder.Methods[Permit2PermitBatch].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	return nil
}

type Permit2PermitTransferFromBatchParams struct {
}

// abi.decode(inputs, (IAllowanceTransfer.AllowanceTransferDetails[]));

func (p *Permit2PermitTransferFromBatchParams) Decode(ctx context.Context, data []byte) error {
	args := make(map[string]interface{})
	err := UniversalRouterDecoder.Methods[Permit2TransferFromBatch].Inputs.UnpackIntoMap(args, data)
	if err != nil {
		return err
	}
	return nil
}
