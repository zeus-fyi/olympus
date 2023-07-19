package artemis_trading_auxiliary

import (
	"context"
	"errors"
	"fmt"

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

func (a *AuxiliaryTradingUtils) UniversalRouterCmdExecutor(ctx context.Context, w3c web3_client.Web3Client, ur *web3_client.UniversalRouterExecCmd) (*types.Transaction, error) {
	signedTx, _, err := universalRouterCmdToTxBuilder(ctx, w3c, ur)
	if err != nil {
		log.Err(err).Msg("error building signed tx")
		return nil, err
	}
	return a.universalRouterExecuteTx(ctx, signedTx)
}

func (a *AuxiliaryTradingUtils) universalRouterExecuteTx(ctx context.Context, signedTx *types.Transaction) (*types.Transaction, error) {
	err := a.w3a().SendSignedTransaction(ctx, signedTx)
	if err != nil {
		log.Err(err).Msg("error sending signed tx")
		return nil, err
	}
	return signedTx, err
}

func debugPrintBalances(ctx context.Context, w3c web3_client.Web3Client) error {
	bal, err := checkEthBalance(ctx, w3c)
	if err != nil {
		return err
	}
	fmt.Println("ETH Balance: ", bal.String())
	bal, err = CheckAuxWETHBalance(ctx, w3c)
	if err != nil {
		return err
	}
	fmt.Println("WETH Balance: ", bal.String())
	return nil
}

// takes a universal router command and returns a signed tx
func universalRouterCmdToTxBuilder(ctx context.Context, w3c web3_client.Web3Client, ur *web3_client.UniversalRouterExecCmd) (*types.Transaction, *web3_actions.SendContractTxPayload, error) {
	if ur == nil {
		return nil, nil, errors.New("universal router command is nil")
	}
	ur.Deadline = GetDeadline()
	data, err := ur.EncodeCommands(ctx)
	if err != nil {
		log.Warn().Msg("universalRouterCmdToTxBuilder: error encoding commands")
		log.Err(err).Msg("error encoding commands")
		return nil, nil, err
	}
	scInfo, err := GetUniswapUniversalRouterAbiPayload(ctx, w3c, data)
	if err != nil {
		log.Warn().Msg("universalRouterCmdToTxBuilder: error getting uniswap universal router abi payload")
		log.Err(err).Msg("error getting uniswap universal router abi payload")
		return nil, nil, err
	}
	signedTx, err := w3c.GetSignedTxToCallFunctionWithData(ctx, &scInfo, scInfo.Data)
	if err != nil {
		log.Warn().Msg("w3c.GetSignedTxToCallFunctionWithData: error getting signed tx to call function with data")
		log.Err(err).Msg("error getting signed tx to call function with data")
		return nil, nil, err
	}
	err = universalRouterCmdVerifier(ctx, w3c, ur, &scInfo)
	if err != nil {
		log.Warn().Msg("universalRouterCmdVerifier: error verifying universal router command")
		log.Err(err).Msg("error verifying universal router command")
		return nil, nil, err
	}
	return signedTx, &scInfo, nil
}

func GetUniswapUniversalRouterAbiPayload(ctx context.Context, w3c web3_client.Web3Client, payload *web3_client.UniversalRouterExecParams) (web3_actions.SendContractTxPayload, error) {
	if payload == nil {
		log.Warn().Msg("GetUniswapUniversalRouterAbiPayload: payload is nil")
		return web3_actions.SendContractTxPayload{}, errors.New("payload is nil")
	}
	payable := payload.Payable
	if payable == nil {
		payable = &web3_actions.SendEtherPayload{
			TransferArgs:   web3_actions.TransferArgs{},
			GasPriceLimits: web3_actions.GasPriceLimits{},
		}
	}
	if payload.Deadline == nil {
		payload.Deadline = GetDeadline()
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
		log.Warn().Msg("GetUniswapUniversalRouterAbiPayload: error generating bin data from params abi")
		log.Err(err).Msg("error generating bin data from params abi")
		return web3_actions.SendContractTxPayload{}, err
	}
	err = w3c.SuggestAndSetGasPriceAndLimitForTx(ctx, &params, common.HexToAddress(params.SmartContractAddr))
	if err != nil {
		log.Warn().Msg("GetUniswapUniversalRouterAbiPayload: error generating bin data from params abi")
		log.Err(err).Msg("error generating bin data from params abi")
		return web3_actions.SendContractTxPayload{}, err
	}
	err = txGasAdjuster(ctx, &params)
	if err != nil {
		log.Warn().Msg("GetUniswapUniversalRouterAbiPayload: error adjusting tx gas")
		log.Err(err).Msg("error adjusting tx gas")
		return web3_actions.SendContractTxPayload{}, err
	}
	return params, nil
}

func checkIfCmdEmpty(ur *web3_client.UniversalRouterExecCmd) *web3_client.UniversalRouterExecCmd {
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
