package artemis_trading_auxiliary

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

func (a *AuxiliaryTradingUtils) addPermit2Ctx(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, TradeCfg, Permit2)
	return ctx
}

func (a *AuxiliaryTradingUtils) SetPermit2ApprovalForToken(ctx context.Context, address string) (*types.Transaction, error) {
	tx, err := a.w3c().ApprovePermit2(ctx, address)
	if err != nil {
		log.Err(err).Msg("error approving permit2")
		return tx, err
	}
	return tx, nil
}

func (a *AuxiliaryTradingUtils) generatePermit2Approval(ctx context.Context, w3c web3_client.Web3Client, to *artemis_trading_types.TradeOutcome) (web3_client.Permit2PermitParams, *artemis_eth_txs.Permit2Tx, error) {
	deadline := GetDeadline()
	chainID, err := getChainID(ctx, w3c)
	if err != nil {
		log.Warn().Err(err).Msg("error getting chainID")
		return web3_client.Permit2PermitParams{}, nil, err
	}
	owner := w3c.Address()
	token := to.AmountInAddr
	spender := artemis_trading_constants.UniswapUniversalRouterNewAddressAccount
	ptNonce, err := GetNextPermit2NonceFromContract(ctx, w3c, owner, token, spender)
	if err != nil {
		return web3_client.Permit2PermitParams{}, nil, err
	}
	pt := &artemis_eth_txs.Permit2Tx{
		Permit2Tx: artemis_autogen_bases.Permit2Tx{
			Owner:             owner.String(),
			Deadline:          int(deadline.Int64()),
			Token:             token.String(),
			ProtocolNetworkID: chainID,
			Nonce:             int(ptNonce.Int64()),
		},
	}
	log.Info().Str("token", token.String()).Str("owner", owner.String()).Int64("nonce", ptNonce.Int64()).Msg("permit2 nonce")
	psp := web3_client.Permit2PermitParams{
		PermitSingle: web3_client.PermitSingle{
			PermitDetails: web3_client.PermitDetails{
				Token:      token,
				Amount:     to.AmountIn,
				Expiration: deadline,
				Nonce:      ptNonce,
			},
			Spender:     spender,
			SigDeadline: deadline,
		},
	}
	err = psp.SignPermit2(w3c.Account, chainID)
	if err != nil {
		log.Warn().Err(err).Msg("error signing permit")
		return web3_client.Permit2PermitParams{}, nil, err
	}
	if psp.Signature == nil {
		log.Warn().Msg("signature is nil")
		return web3_client.Permit2PermitParams{}, nil, err
	}
	return psp, pt, err
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

func GetNextPermit2NonceFromContract(ctx context.Context, w3c web3_client.Web3Client, owner, token, spender accounts.Address) (*big.Int, error) {
	nonceCheck := &web3_actions.SendContractTxPayload{
		SmartContractAddr: artemis_trading_constants.Permit2SmartContractAddress,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       artemis_oly_contract_abis.MustLoadPermit2Abi(),
		MethodName:        "allowance",
		Params:            []interface{}{owner, token, spender},
	}
	resp, err := w3c.CallConstantFunction(ctx, nonceCheck)
	if err != nil {
		log.Err(err).Msg("error getting nonce")
		return nil, err
	}
	if len(resp) != 3 {
		return nil, errors.New("unexpected response")
	}
	for i, val := range resp {
		switch i {
		case 2:
			return val.(*big.Int), nil
		}
	}
	return nil, nil
}
