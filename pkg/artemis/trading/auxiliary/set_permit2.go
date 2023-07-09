package artemis_trading_auxiliary

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

func (a *AuxiliaryTradingUtils) SetPermit2ApprovalForToken(ctx context.Context, address string) (*types.Transaction, error) {
	tx, err := a.ApprovePermit2(ctx, address)
	if err != nil {
		log.Err(err).Msg("error approving permit2")
		return tx, err
	}
	return tx, nil
}

// AuxiliaryTradingUtils GetNonce: todo this needs to update a nonce count in db or track them somehow

var ts chronos.Chronos

func (a *AuxiliaryTradingUtils) GetPermit2Nonce() *big.Int {
	nonce := new(big.Int).SetUint64(uint64(ts.GeneratePermit2Nonce()))
	return nonce
}

func (a *AuxiliaryTradingUtils) generatePermit2Approval(ctx context.Context, to *artemis_trading_types.TradeOutcome) (web3_client.Permit2PermitParams, error) {
	deadline := a.GetDeadline()
	psp := web3_client.Permit2PermitParams{
		PermitSingle: web3_client.PermitSingle{
			PermitDetails: web3_client.PermitDetails{
				Token:      to.AmountInAddr,
				Amount:     to.AmountIn,
				Expiration: deadline,
				Nonce:      a.GetPermit2Nonce(),
			},
			Spender:     artemis_trading_constants.UniswapUniversalRouterNewAddressAccount,
			SigDeadline: deadline,
		},
	}
	chainID, err := a.getChainID(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("error getting chainID")
		return psp, err
	}
	err = psp.SignPermit2(a.Account, chainID)
	if err != nil {
		log.Warn().Err(err).Msg("error signing permit")
		return psp, err
	}
	if psp.Signature == nil {
		log.Warn().Msg("signature is nil")
		return psp, errors.New("signature is nil")
	}
	return psp, err
}

func (a *AuxiliaryTradingUtils) generatePermit2Transfer(ctx context.Context, to *artemis_trading_types.TradeOutcome) (web3_client.Permit2TransferFromParams, error) {
	deadline := a.GetDeadline()
	psp := web3_client.Permit2TransferFromParams{
		PermitTransferFrom: web3_client.PermitTransferFrom{
			TokenPermissions: web3_client.TokenPermissions{
				Token:  to.AmountInAddr,
				Amount: to.AmountIn,
			},
			Nonce:    a.GetPermit2Nonce(),
			Deadline: deadline,
		},
		Permit2SignatureTransferDetails: web3_client.Permit2SignatureTransferDetails{
			To:              artemis_trading_constants.UniswapUniversalRouterNewAddressAccount,
			RequestedAmount: to.AmountIn,
		},
		Owner:     a.Address(),
		Signature: nil,
	}
	chainID, err := a.getChainID(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("error getting chainID")
		return psp, err
	}
	err = psp.SignPermit2(a.Account, chainID)
	if err != nil {
		log.Warn().Err(err).Msg("error signing permit")
		return psp, err
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
