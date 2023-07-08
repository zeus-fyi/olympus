package artemis_realtime_trading

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

// TODO add v3 support

func (a *ActiveTrading) ExecTradeV2SwapFromTokenToToken(ctx context.Context, to *artemis_trading_types.TradeOutcome, bypassSim bool) error {
	// todo max this window more appropriate vs near infinite
	sigDeadline, _ := new(big.Int).SetString("3000000000000", 10)
	ur := web3_client.UniversalRouterExecCmd{
		Commands: []web3_client.UniversalRouterExecSubCmd{},
		Deadline: sigDeadline,
		Payable:  nil,
	}
	sc1 := web3_client.UniversalRouterExecSubCmd{
		Command:   artemis_trading_constants.Permit2Permit,
		CanRevert: false,
		Inputs:    nil,
	}
	psp := web3_client.Permit2PermitParams{
		PermitSingle: web3_client.PermitSingle{
			PermitDetails: web3_client.PermitDetails{
				Token:      to.AmountInAddr,
				Amount:     to.AmountIn,
				Expiration: sigDeadline,
				// todo this needs to update a nonce count in db or track them somehow
				Nonce: new(big.Int).SetUint64(0),
			},
			Spender:     artemis_trading_constants.UniswapUniversalRouterNewAddressAccount,
			SigDeadline: sigDeadline,
		},
	}
	err := psp.SignPermit2Mainnet(a.u.Web3Client.Account)
	if err != nil {
		log.Warn().Err(err).Msg("error signing permit")
		return err
	}
	if psp.Signature == nil {
		log.Warn().Msg("signature is nil")
		return errors.New("signature is nil")
	}
	sc1.DecodedInputs = psp
	ur.Commands = append(ur.Commands, sc1)
	sc2 := web3_client.UniversalRouterExecSubCmd{
		Command:   artemis_trading_constants.V2SwapExactIn,
		CanRevert: false,
		Inputs:    nil,
		DecodedInputs: web3_client.V2SwapExactInParams{
			AmountIn:      to.AmountIn,
			AmountOutMin:  to.AmountOut,
			Path:          []accounts.Address{to.AmountInAddr, to.AmountOutAddr},
			To:            accounts.HexToAddress(artemis_trading_constants.UniversalRouterSender),
			PayerIsSender: true,
		},
	}
	ur.Commands = append(ur.Commands, sc2)
	tx, err := a.execUniswapUniversalRouterCmd(ctx, ur, bypassSim)
	if err != nil {
		return err
	}
	to.BundleTxs = append(to.BundleTxs, tx)
	to.AddTxHash(accounts.Hash(tx.Hash()))
	return err
}

func (a *ActiveTrading) execUniswapUniversalRouterCmd(ctx context.Context, payload web3_client.UniversalRouterExecCmd, bypassSim bool) (*types.Transaction, error) {
	data, err := payload.EncodeCommands(ctx)
	if err != nil {
		log.Err(err).Msg("ExecUniswapUniversalRouterCmd: failed to encode commands")
		return nil, err
	}
	scInfo := GetUniswapUniversalRouterAbiPayload(data)
	signedTx, err := a.u.Web3Client.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
	if err != nil {
		return nil, err
	}
	if bypassSim {
		return signedTx, nil
	}
	err = a.u.Web3Client.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func GetUniswapUniversalRouterAbiPayload(payload *web3_client.UniversalRouterExecParams) web3_actions.SendContractTxPayload {
	if payload == nil {
		payload = &web3_client.UniversalRouterExecParams{}
		log.Warn().Msg("GetUniswapUniversalRouterAbiPayload: payload is nil")
		return web3_actions.SendContractTxPayload{}
	}
	payable := payload.Payable
	if payable == nil {
		payable = &web3_actions.SendEtherPayload{
			TransferArgs:   web3_actions.TransferArgs{},
			GasPriceLimits: web3_actions.GasPriceLimits{},
		}
	}
	fnParams := []interface{}{payload.Commands, payload.Inputs}
	methodName := artemis_trading_constants.Execute
	if payload.Deadline != nil {
		methodName = artemis_trading_constants.Execute0
		fnParams = append(fnParams, payload.Deadline.String())
	}
	params := web3_actions.SendContractTxPayload{
		SmartContractAddr: artemis_trading_constants.UniswapUniversalRouterAddressNew,
		SendEtherPayload:  *payable,
		ContractABI:       artemis_trading_constants.UniversalRouterNewAbi,
		MethodName:        methodName,
		Params:            fnParams,
	}
	return params
}
