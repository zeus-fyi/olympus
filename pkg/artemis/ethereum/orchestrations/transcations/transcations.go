package artemis_ethereum_transcations

import (
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/aegis/web3_client"
)

var ArtemisEthereumBroadcastTxClient web3_client.Web3Client

func InitArtemisEthereumClient(nodeURL string, acc *accounts.Account) {
	ArtemisEthereumBroadcastTxClient = web3_client.NewWeb3Client(nodeURL, acc)
}
