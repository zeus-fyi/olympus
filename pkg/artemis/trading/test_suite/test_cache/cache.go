package artemis_test_cache

import (
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
)

var LiveTestNetwork = web3_actions.Web3Actions{}

func InitLiveTestNetwork(nodeURL string) {
	LiveTestNetwork = web3_actions.NewWeb3ActionsClient(nodeURL)
	LiveTestNetwork.AddDefaultEthereumMainnetTableHeader()
}
