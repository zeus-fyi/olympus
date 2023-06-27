package metrics_trading

import "github.com/zeus-fyi/olympus/pkg/artemis/web3_client"

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
	tx.Stats.WithLabelValues(label, method).Inc()
}
