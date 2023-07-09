package artemis_trading_auxiliary

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) SetPermit2ApprovalForToken(ctx context.Context, address string) (*types.Transaction, error) {
	tx, err := a.ApprovePermit2(ctx, address)
	if err != nil {
		log.Err(err).Msg("error approving permit2")
		return tx, err
	}
	return tx, nil
}

func (a *AuxiliaryTradingUtils) GetDeadline() *big.Int {
	deadline := int(time.Now().Add(60 * time.Second).Unix())
	sigDeadline := artemis_eth_units.NewBigInt(deadline)
	return sigDeadline
}

// AuxiliaryTradingUtils GetNonce: todo this needs to update a nonce count in db or track them somehow

func (a *AuxiliaryTradingUtils) GetPermit2Nonce() *big.Int {
	nonce := new(big.Int).SetUint64(0)
	return nonce
}

func (a *AuxiliaryTradingUtils) generatePermit2Approval(ctx context.Context, tokenAddr accounts.Address, amount *big.Int) (web3_client.Permit2PermitParams, error) {
	deadline := a.GetDeadline()
	psp := web3_client.Permit2PermitParams{
		PermitSingle: web3_client.PermitSingle{
			PermitDetails: web3_client.PermitDetails{
				Token:      tokenAddr,
				Amount:     amount,
				Expiration: deadline,
				Nonce:      a.GetPermit2Nonce(),
			},
			Spender:     artemis_trading_constants.UniswapUniversalRouterNewAddressAccount,
			SigDeadline: deadline,
		},
	}
	err := psp.SignPermit2Mainnet(a.Account)
	if err != nil {
		log.Err(err).Msg("error signing permit")
		return web3_client.Permit2PermitParams{}, err
	}
	if psp.Signature == nil {
		log.Err(err).Msg("signature is nil")
		return web3_client.Permit2PermitParams{}, errors.New("signature is nil")
	}
	return psp, err
}

func (a *AuxiliaryTradingUtils) generatePermit2Transfer(ctx context.Context, tokenAddr accounts.Address, amount *big.Int) (web3_client.Permit2TransferFromParams, error) {
	deadline := a.GetDeadline()
	psp := web3_client.Permit2TransferFromParams{
		PermitTransferFrom: web3_client.PermitTransferFrom{
			TokenPermissions: web3_client.TokenPermissions{
				Token:  tokenAddr,
				Amount: amount,
			},
			Nonce:    a.GetPermit2Nonce(),
			Deadline: deadline,
		},
		Permit2SignatureTransferDetails: web3_client.Permit2SignatureTransferDetails{
			To:              artemis_trading_constants.UniswapUniversalRouterNewAddressAccount,
			RequestedAmount: amount,
		},
		Owner:     a.Address(),
		Signature: nil,
	}
	err := psp.SignPermit2Mainnet(a.Account)
	if err != nil {
		log.Err(err).Msg("error signing permit")
		return web3_client.Permit2TransferFromParams{}, err
	}
	if psp.Signature == nil {
		log.Err(err).Msg("signature is nil")
		return web3_client.Permit2TransferFromParams{}, errors.New("signature is nil")
	}
	return psp, err
}

/*
	permit, err := a.generatePermit2Transfer(ctx, wethAddr, amountIn)
	if err != nil {
		return nil, err
	}
	permitCmd := web3_client.UniversalRouterExecSubCmd{
		Command:       artemis_trading_constants.Permit2TransferFrom,
		DecodedInputs: permit,
		CanRevert:     false,
	}
	ur.Commands = append(ur.Commands, permitCmd)
*/
// todo set batch here
