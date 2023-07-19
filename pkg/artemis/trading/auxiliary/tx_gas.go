package artemis_trading_auxiliary

import (
	"context"

	"github.com/rs/zerolog/log"
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

tx.GasFeeCap()) 14009230241
tx.GasTipCap() 0
tx.GasPrice() 14009230241
tx.Gas() 165060

tx.GasFeeCap()) 17216635871
tx.GasTipCap() 100000000
tx.GasPrice() 17216635871
tx.Gas() 326436

tx.GasFeeCap()) 36180761500
tx.GasTipCap() 36180761500
tx.GasPrice() 36180761500
tx.Gas() 142255
*/

func txGasAdjuster(ctx context.Context, scInfo *web3_actions.SendContractTxPayload) error {
	tt := getTradeTypeFromCtx(ctx)
	switch tt {
	case FrontRun:
		log.Info().Msg("txGasAdjuster: FrontRun gas adjustment")
		scInfo.GasTipCap = artemis_eth_units.NewBigInt(1)
	case UserTrade:
		log.Info().Msg("txGasAdjuster: UserTrade gas adjustment")
		scInfo.GasTipCap = artemis_eth_units.NewBigInt(1)
		scInfo.GasLimit *= 2
	case BackRun:
		log.Info().Msg("txGasAdjuster: BackRun gas adjustment")
		scInfo.GasFeeCap = artemis_eth_units.MulBigIntFromInt(scInfo.GasFeeCap, 2)
		scInfo.GasTipCap = artemis_eth_units.MulBigIntFromInt(scInfo.GasTipCap, 2)
	default:
		return nil
	}
	return nil
}
