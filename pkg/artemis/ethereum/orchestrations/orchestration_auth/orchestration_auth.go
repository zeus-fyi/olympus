package artemis_orchestration_auth

import (
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	mempool_txs "github.com/zeus-fyi/olympus/datastores/dynamodb/mempool"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

var (
	Bearer            string
	MevDynamoDBClient mempool_txs.MempoolTxDynamoDB
)

func InitMevDynamoDBClient(creds dynamodb_client.DynamoDBCredentials) {
	wc := web3_client.NewWeb3Client(artemis_network_cfgs.ArtemisEthereumMainnetQuiknode.NodeURL, artemis_network_cfgs.ArtemisEthereumMainnet.Account)
	MevDynamoDBClient = mempool_txs.NewMempoolTxDynamoDB(creds)
}
