package artemis_trading_auxiliary

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

func (a *AuxiliaryTradingUtils) UniversalRouterCmdExecutor(ctx context.Context, ur *web3_client.UniversalRouterExecCmd) (*types.Transaction, error) {
	signedTx, err := a.universalRouterCmdBuilder(ctx, ur)
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
	a.addTx(signedTx)
	return signedTx, err
}

/*
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64 // a.k.a. gasLimit
*/

func (a *AuxiliaryTradingUtils) txGasAdjuster(ctx context.Context, tx *types.Transaction, cmdType string) (*types.Transaction, error) {
	scInfo := &web3_actions.SendContractTxPayload{
		SmartContractAddr: "",
		SendEtherPayload: web3_actions.SendEtherPayload{
			TransferArgs: web3_actions.TransferArgs{
				Amount: tx.Value(),
			},
			GasPriceLimits: web3_actions.GasPriceLimits{
				GasPrice:  tx.GasPrice(),
				GasLimit:  tx.Gas(),
				GasTipCap: tx.GasTipCap(),
				GasFeeCap: tx.GasFeeCap(),
			},
		},
	}
	err := a.SuggestAndSetGasPriceAndLimitForTx(ctx, scInfo, common.HexToAddress(scInfo.ToAddress.Hex()), tx.Data())
	if err != nil {
		return nil, err
	}
	switch cmdType {
	case "frontRun":
		scInfo.GasTipCap = artemis_eth_units.Finney
	case "sandwich":
		scInfo.GasTipCap = artemis_eth_units.GweiFraction(1, 10)
	case "backRun":
		scInfo.GasTipCap = artemis_eth_units.MulBigIntFromInt(scInfo.GasTipCap, 2)
	}
	jtx := artemis_trading_types.JSONTx{}
	err = jtx.UnmarshalTx(tx)
	if err != nil {
		return nil, err
	}
	//cid, err := a.getChainID(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//jtx.ChainID = strconv.Itoa(cid)
	//jtx.Gas = strconv.FormatUint(scInfo.GasLimit, 10)
	//jtx.MaxFeePerGas = scInfo.GasFeeCap
	jtx.MaxPriorityFeePerGas = scInfo.GasTipCap
	tx, err = jtx.ConvertToTx()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (a *AuxiliaryTradingUtils) universalRouterCmdBuilder(ctx context.Context, ur *web3_client.UniversalRouterExecCmd) (*types.Transaction, error) {
	ur.Deadline = a.GetDeadline()
	data, err := ur.EncodeCommands(ctx)
	if err != nil {
		return nil, err
	}
	scInfo := GetUniswapUniversalRouterAbiPayload(data)
	signedTx, err := a.GetSignedTxToCallFunctionWithArgs(ctx, &scInfo)
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

var urAbi = artemis_oly_contract_abis.MustLoadNewUniversalRouterAbi()

func GetUniswapUniversalRouterAbiPayload(payload *web3_client.UniversalRouterExecParams) web3_actions.SendContractTxPayload {
	if payload == nil {
		payload = &web3_client.UniversalRouterExecParams{}
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
		ContractABI:       urAbi,
		MethodName:        methodName,
		Params:            fnParams,
	}
	return params
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

func (a *AuxiliaryTradingUtils) GetDeadline() *big.Int {
	deadline := int(time.Now().Add(60 * time.Second).Unix())
	sigDeadline := artemis_eth_units.NewBigInt(deadline)
	return sigDeadline
}
