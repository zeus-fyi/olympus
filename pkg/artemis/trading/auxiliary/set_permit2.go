package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (a *AuxiliaryTradingUtils) addPermit2Ctx(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, TradeCfg, Permit2)
	return ctx
}

func (a *AuxiliaryTradingUtils) SetPermit2ApprovalForToken(ctx context.Context, address string) (*types.Transaction, error) {
	tx, err := a.getWeb3Client().ApprovePermit2(ctx, address)
	if err != nil {
		log.Err(err).Msg("error approving permit2")
		return tx, err
	}
	return tx, nil
}

func (a *AuxiliaryTradingUtils) generatePermit2Approval(ctx context.Context, to *artemis_trading_types.TradeOutcome) (web3_client.Permit2PermitParams, error) {
	deadline := a.GetDeadline()
	chainID, err := a.getChainID(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("error getting chainID")
		return web3_client.Permit2PermitParams{}, err
	}
	pt := &artemis_eth_txs.Permit2Tx{
		Permit2Tx: artemis_autogen_bases.Permit2Tx{
			Owner:             a.U.Web3Client.Address().String(),
			Deadline:          int(deadline.Int64()),
			Token:             to.AmountInAddr.String(),
			ProtocolNetworkID: chainID,
		},
	}
	err = pt.SelectNextPermit2Nonce(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("error getting permit2 nonce")
		return web3_client.Permit2PermitParams{}, err
	}
	ptNonce := artemis_eth_units.NewBigInt(pt.Nonce)
	psp := web3_client.Permit2PermitParams{
		PermitSingle: web3_client.PermitSingle{
			PermitDetails: web3_client.PermitDetails{
				Token:      to.AmountInAddr,
				Amount:     to.AmountIn,
				Expiration: deadline,
				Nonce:      ptNonce,
			},
			Spender:     artemis_trading_constants.UniswapUniversalRouterNewAddressAccount,
			SigDeadline: deadline,
		},
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
	panic("not implemented")
	//deadline := a.GetDeadline()
	//psp := web3_client.Permit2TransferFromParams{
	//	PermitTransferFrom: web3_client.PermitTransferFrom{
	//		TokenPermissions: web3_client.TokenPermissions{
	//			Token:  to.AmountInAddr,
	//			Amount: to.AmountIn,
	//		},
	//		Nonce:    a.GetPermit2Nonce(),
	//		Deadline: deadline,
	//	},
	//	Permit2SignatureTransferDetails: web3_client.Permit2SignatureTransferDetails{
	//		To:              artemis_trading_constants.UniswapUniversalRouterNewAddressAccount,
	//		RequestedAmount: to.AmountIn,
	//	},
	//	Owner:     a.Address(),
	//	Signature: nil,
	//}
	//chainID, err := a.getChainID(ctx)
	//if err != nil {
	//	log.Warn().Err(err).Msg("error getting chainID")
	//	return psp, err
	//}
	//err = psp.SignPermit2(a.Account, chainID)
	//if err != nil {
	//	log.Warn().Err(err).Msg("error signing permit")
	//	return psp, err
	//}
	//return psp, err
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
