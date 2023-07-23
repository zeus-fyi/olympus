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
{"level":"info","builder":"https://api.blocknative.com/v1/auction","bundleHash":"0x60d198f17e43c0ed4c0a4f2a7a39b4b89014be20d63076e1fdac47a6056223a8","time":1690073560,"message":"sendAdditionalBundles: bundle sent successfully"}
{"level":"info","builder":"https://rpc.nfactorial.xyz/","bundleHash":"0x60d198f17e43c0ed4c0a4f2a7a39b4b89014be20d63076e1fdac47a6056223a8","time":1690073560,"message":"sendAdditionalBundles: bundle sent successfully"}
{"level":"info","resp":{"bundleGasPrice":"4022970258","bundleHash":"0x60d198f17e43c0ed4c0a4f2a7a39b4b89014be20d63076e1fdac47a6056223a8","coinbaseDiff":"967427795834944","ethSentToCoinbase":"0","gasFees":"967427795834944",
"results":[{"coinbaseDiff":"0","ethSentToCoinbase":"0","fromAddress":"0x000000641e80A183c8B736141cbE313E136bc8c6","gasFees":"0","gasPrice":"0","gasUsed":32611,"toAddress":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
"txHash":"0x3555daa7f8a5e81effbf7d81f209a67c29a184298a9453390a3149a09e19e891","value":"","error":"out of gas","revert":""},
{"coinbaseDiff":"17199300000000","ethSentToCoinbase":"0","fromAddress":"0x0F012414F3E774C29f0d4b1311B8f87Ba396fE6C","gasFees":"17199300000000","gasPrice":"100000000","gasUsed":171993,"toAddress":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
"txHash":"0x79b78f96d34be9388b32570072f5788138bf63dad4773c43b068cc114b57c5a2","value":"0x","error":"","revert":""},
{"coinbaseDiff":"950228495834944","ethSentToCoinbase":"0","fromAddress":"0x000000641e80A183c8B736141cbE313E136bc8c6","gasFees":"950228495834944","gasPrice":"26489420602","gasUsed":35872,
"toAddress":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","txHash":"0x734b8334d18ac36a1e5e8b0607c63923ac88d3c01e7c4de592547ea6ceb27eb0","value":"","error":"out of gas","revert":""}],
"stateBlockNumber":17752361,"totalGasUsed":240476},"resp.BundleGasPrice":"4022970258","fbCallResp":[{"coinbaseDiff":"0","ethSentToCoinbase":"0","fromAddress":"0x000000641e80A183c8B736141cbE313E136bc8c6"
,"gasFees":"0","gasPrice":"0","gasUsed":32611,"toAddress":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","txHash":"0x3555daa7f8a5e81effbf7d81f209a67c29a184298a9453390a3149a09e19e891","value":"","error":"out of gas","revert":""},
{"coinbaseDiff":"17199300000000","ethSentToCoinbase":"0","fromAddress":"0x0F012414F3E774C29f0d4b1311B8f87Ba396fE6C","gasFees":"17199300000000","gasPrice":"100000000","gasUsed":171993,"toAddress":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD",
"txHash":"0x79b78f96d34be9388b32570072f5788138bf63dad4773c43b068cc114b57c5a2","value":"0x","error":"","revert":""},{"coinbaseDiff":"950228495834944","ethSentToCoinbase":"0","fromAddress":"0x000000641e80A183c8B736141cbE313E136bc8c6","gasFees":"950228495834944","gasPrice":"26489420602","gasUsed":35872,
"toAddress":"0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD","txHash":"0x734b8334d18ac36a1e5e8b0607c63923ac88d3c01e7c4de592547ea6ceb27eb0","value":"","error":"out of gas","revert":""}],"time":1690073561,"message":"CallFlashbotsBundle: bundle sent successfully"}
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
