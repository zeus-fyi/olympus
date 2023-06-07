package web3_client

import (
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
	Permit2Permit            = "PERMIT2_PERMIT"
	Permit2TransferFrom      = "PERMIT2_TRANSFER_FROM"
	Permit2PermitBatch       = "PERMIT2_PERMIT_BATCH"
	Permit2TransferFromBatch = "PERMIT2_TRANSFER_FROM_BATCH"
)

type Permit2PermitTransferFrom struct {
	TokenAddr accounts.Address `json:"tokenAddr"`
	Recipient accounts.Address `json:"recipient"`
	Amount    *big.Int         `json:"amount"`
}

type Permit2PermitVal struct {
	Permit2PermitTransferFrom
	Signature []byte `json:"signature"`
}
