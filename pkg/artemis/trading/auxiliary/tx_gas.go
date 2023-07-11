package artemis_trading_auxiliary

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
)

/*
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64 // a.k.a. gasLimit
*/

func (a *AuxiliaryTradingUtils) txGasAdjuster(ctx context.Context, txWithMetadata TxWithMetadata) (*types.Transaction, error) {
	tx := txWithMetadata.Tx
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
	switch txWithMetadata.TradeType {
	case "frontRun":
		scInfo.GasTipCap = artemis_eth_units.Finney
	case "sandwich":
		scInfo.GasTipCap = artemis_eth_units.GweiFraction(1, 10)
	case "backRun":
		scInfo.GasTipCap = artemis_eth_units.MulBigIntFromInt(scInfo.GasTipCap, 2)
	default:
		return tx, nil
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
