package artemis_ethereum_transcations

import web3_client "github.com/zeus-fyi/olympus/pkg/aegis/web3_client"

var ArtemisEthereumBroadcastTxClient web3_client.Web3Client

func InitArtemisEthereumClient(nodeURL string) {
	ArtemisEthereumBroadcastTxClient = web3_client.NewClient(nodeURL)
}
