package artemis_trading_auxiliary

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

var (
	urAbi = artemis_oly_contract_abis.MustLoadNewUniversalRouterAbi()
)

func (a *AuxiliaryTradingUtils) UniversalRouterCmdExecutor(ctx context.Context, ur *web3_client.UniversalRouterExecCmd) (*types.Transaction, error) {
	signedTx, err := a.universalRouterCmdToTxBuilder(ctx, ur)
	if err != nil {
		log.Err(err).Msg("error building signed tx")
		return nil, err
	}
	return a.universalRouterExecuteTx(ctx, signedTx)
}

func (a *AuxiliaryTradingUtils) universalRouterExecuteTx(ctx context.Context, signedTx *types.Transaction) (*types.Transaction, error) {
	err := a.SendSignedTransaction(ctx, signedTx)
	if err != nil {
		log.Err(err).Msg("error sending signed tx")
		return nil, err
	}
	return signedTx, err
}

func (a *AuxiliaryTradingUtils) universalRouterCmdToTxBuilder(ctx context.Context, ur *web3_client.UniversalRouterExecCmd) (*types.Transaction, error) {
	ur.Deadline = a.GetDeadline()
	data, err := ur.EncodeCommands(ctx)
	if err != nil {
		return nil, err
	}
	scInfo, err := a.GetUniswapUniversalRouterAbiPayload(ctx, data)
	if err != nil {
		return nil, err
	}
	signedTx, err := a.GetSignedTxToCallFunctionWithData(ctx, &scInfo, scInfo.Data)
	if err != nil {
		return signedTx, err
	}
	err = a.universalRouterCmdVerifier(ctx, ur, &scInfo)
	if err != nil {
		return nil, err
	}
	err = a.AddTxToBundleGroup(ctx, signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (a *AuxiliaryTradingUtils) GetUniswapUniversalRouterAbiPayload(ctx context.Context, payload *web3_client.UniversalRouterExecParams) (web3_actions.SendContractTxPayload, error) {
	if payload == nil {
		payload = &web3_client.UniversalRouterExecParams{}
		return web3_actions.SendContractTxPayload{}, errors.New("payload is nil")
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
		ContractABI:       urAbi,
		MethodName:        methodName,
		Params:            fnParams,
	}
	err := params.GenerateBinDataFromParamsAbi(ctx)
	if err != nil {
		log.Err(err).Msg("error generating bin data from params abi")
		return web3_actions.SendContractTxPayload{}, err
	}
	err = a.SuggestAndSetGasPriceAndLimitForTx(ctx, &params, common.HexToAddress(params.SmartContractAddr))
	if err != nil {
		return web3_actions.SendContractTxPayload{}, err
	}
	err = a.txGasAdjuster(ctx, &params)
	if err != nil {
		return web3_actions.SendContractTxPayload{}, err
	}
	return params, nil
}

func (a *AuxiliaryTradingUtils) checkIfCmdEmpty(ur *web3_client.UniversalRouterExecCmd) *web3_client.UniversalRouterExecCmd {
	if ur == nil {
		ur = &web3_client.UniversalRouterExecCmd{
			Commands: []web3_client.UniversalRouterExecSubCmd{},
			Payable: &web3_actions.SendEtherPayload{
				TransferArgs: web3_actions.TransferArgs{
					Amount:    nil,
					ToAddress: accounts.Address{},
				},
				GasPriceLimits: web3_actions.GasPriceLimits{
					GasPrice:  nil,
					GasLimit:  0,
					GasTipCap: nil,
					GasFeeCap: nil,
				},
			},
		}
	}
	if ur.Commands == nil {
		ur.Commands = []web3_client.UniversalRouterExecSubCmd{}
	}
	return ur
}
