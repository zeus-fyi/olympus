package artemis_trading_auxiliary

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type AuxiliaryTradingUtils struct {
	web3_client.Web3Client
}

func (a *AuxiliaryTradingUtils) SetPermit2Approval(ctx context.Context, address string) (*types.Transaction, error) {
	tx, err := a.ApprovePermit2(ctx, address)
	if err != nil {
		log.Err(err).Msg("error approving permit")
		return tx, err
	}
	return tx, nil
}

// todo, set this based on block time

func (a *AuxiliaryTradingUtils) GetDeadline() *big.Int {
	sigDeadline, _ := new(big.Int).SetString("3000000000000", 10)
	return sigDeadline
}

// AuxiliaryTradingUtils GetNonce: todo this needs to update a nonce count in db or track them somehow

func (a *AuxiliaryTradingUtils) GetNonce() *big.Int {
	nonce := new(big.Int).SetUint64(0)
	return nonce
}

func (a *AuxiliaryTradingUtils) GeneratePermit2Approval(ctx context.Context, tokenAddr accounts.Address, amount *big.Int) (web3_client.Permit2PermitParams, error) {
	psp := web3_client.Permit2PermitParams{
		PermitSingle: web3_client.PermitSingle{
			PermitDetails: web3_client.PermitDetails{
				Token:      tokenAddr,
				Amount:     amount,
				Expiration: a.GetDeadline(),
				Nonce:      a.GetNonce(),
			},
			Spender:     artemis_trading_constants.UniswapUniversalRouterNewAddressAccount,
			SigDeadline: a.GetDeadline(),
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
