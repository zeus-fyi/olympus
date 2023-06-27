package metrics_txfetcher

import "github.com/zeus-fyi/olympus/pkg/artemis/web3_client"

//	UniswapUniversalRouterAddressNew = "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
//	UniswapUniversalRouterAddressOld = "0xEf1c6E67703c7BD7107eed8303Fbe6EC2554BF6B"
//	UniswapV2Router02Address         = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
//	UniswapV2Router01Address         = "0xf164fC0Ec4E93095b804a4795bBe1e041497b92a"
//
//	UniswapV3Router01Address = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
//	UniswapV3Router02Address = "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45"

var AddressLabelMap = map[string]string{
	web3_client.UniswapUniversalRouterAddressNew: "UniswapUniversalRouterAddressNew",
	web3_client.UniswapUniversalRouterAddressOld: "UniswapUniversalRouterAddressOld",
	web3_client.UniswapV2Router02Address:         "UniswapV2Router02Address",
	web3_client.UniswapV2Router01Address:         "UniswapV2Router01Address",
	web3_client.UniswapV3Router01Address:         "UniswapV3Router01Address",
	web3_client.UniswapV3Router02Address:         "UniswapV3Router02Address",
}

func (tx *TxFetcherMetrics) TransactionGroup(address string, method string) {
	label, ok := AddressLabelMap[address]
	if !ok {
		return
	}
	tx.MevTxStats.WithLabelValues(label, method).Inc()
}
