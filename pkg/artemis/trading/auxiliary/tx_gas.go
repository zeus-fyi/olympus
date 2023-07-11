package artemis_trading_auxiliary

import (
	"context"

	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
)

/*
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64 // a.k.a. gasLimit

max priority fee per gas higher than max fee per gas:
address 0x000025e60C7ff32a3470be7FE3ed1666b0E326e2, maxPriorityFeePerGas: 329436, maxFeePerGas: 164742;
txhash 0x5d6466f6026e0fb7b1cada8e52da091215cb3ea322cf650bd3b094012c2df5e1"}
*/

func (a *AuxiliaryTradingUtils) txGasAdjuster(ctx context.Context, scInfo *web3_actions.SendContractTxPayload) error {
	tt := a.getTradeTypeFromCtx(ctx)
	switch tt {
	case FrontRun:
		scInfo.GasTipCap = artemis_eth_units.Finney
	case UserTrade:
		scInfo.GasTipCap = artemis_eth_units.GweiFraction(1, 10)
	case BackRun:
		scInfo.GasTipCap = artemis_eth_units.MulBigIntFromInt(scInfo.GasTipCap, 2)
	default:
		return nil
	}
	return nil
}
