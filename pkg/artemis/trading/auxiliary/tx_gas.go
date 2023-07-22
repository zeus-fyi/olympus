package artemis_trading_auxiliary

import (
	"github.com/ethereum/go-ethereum/core/types"
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

/*
	scInfoSand.GasLimit = frScInfo.GasLimit * 2
	scInfoSand.GasTipCap = artemis_eth_units.MulBigIntWithInt(frScInfo.GasFeeCap, 4)
	scInfoSand.GasFeeCap = artemis_eth_units.MulBigIntWithInt(frScInfo.GasFeeCap, 4)
*/

func ApplyFrontRunGasAdjustment(signedFrontRunTx *types.Transaction) (web3_actions.GasPriceLimits, web3_actions.GasPriceLimits) {
	if signedFrontRunTx == nil {
		return web3_actions.GasPriceLimits{}, web3_actions.GasPriceLimits{}
	}
	gasFeeCap := signedFrontRunTx.GasFeeCap()
	gasTipCap := signedFrontRunTx.GasTipCap()
	gasLimit := signedFrontRunTx.Gas()
	startingGp := web3_actions.GasPriceLimits{
		GasPrice:  nil,
		GasLimit:  gasLimit,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
	}
	newGasTipCap := artemis_eth_units.MulBigIntWithInt(gasFeeCap, 0)
	adjustedGp := web3_actions.GasPriceLimits{
		GasPrice:  nil,
		GasLimit:  gasLimit,
		GasTipCap: newGasTipCap,
		GasFeeCap: gasFeeCap,
	}
	return startingGp, adjustedGp
}

func ApplyBackrunGasAdjustment(signedFrontRunTx *types.Transaction) (web3_actions.GasPriceLimits, web3_actions.GasPriceLimits) {
	if signedFrontRunTx == nil {
		return web3_actions.GasPriceLimits{}, web3_actions.GasPriceLimits{}
	}
	gasFeeCap := signedFrontRunTx.GasFeeCap()
	gasTipCap := signedFrontRunTx.GasTipCap()
	gasLimit := signedFrontRunTx.Gas()
	startingGp := web3_actions.GasPriceLimits{
		GasPrice:  nil,
		GasLimit:  gasLimit,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
	}
	newGasTipCap := artemis_eth_units.MulBigIntWithInt(gasFeeCap, 4)
	newGasFeeCap := artemis_eth_units.MulBigIntWithInt(gasFeeCap, 4)
	newGasLimit := gasLimit * 2
	adjustedGp := web3_actions.GasPriceLimits{
		GasPrice:  nil,
		GasLimit:  newGasLimit,
		GasTipCap: newGasTipCap,
		GasFeeCap: newGasFeeCap,
	}
	return startingGp, adjustedGp
}

func ApplyTxType2UserGasAdjustment(signedTx *types.Transaction) (web3_actions.GasPriceLimits, web3_actions.GasPriceLimits) {
	if signedTx == nil {
		return web3_actions.GasPriceLimits{}, web3_actions.GasPriceLimits{}
	}
	gasFeeCap := signedTx.GasFeeCap()
	gasTipCap := signedTx.GasTipCap()
	gasLimit := signedTx.Gas()

	startingGp := web3_actions.GasPriceLimits{
		GasPrice:  nil,
		GasLimit:  gasLimit,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
	}
	newGasTipCap := artemis_eth_units.OneTenthGwei
	newGasLimit := gasLimit * 2
	adjustedGp := web3_actions.GasPriceLimits{
		GasPrice:  nil,
		GasLimit:  newGasLimit,
		GasTipCap: newGasTipCap,
		GasFeeCap: gasFeeCap,
	}
	return startingGp, adjustedGp
}
