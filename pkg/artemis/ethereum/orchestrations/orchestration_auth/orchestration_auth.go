package artemis_orchestration_auth

import (
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	mempool_txs "github.com/zeus-fyi/olympus/datastores/dynamodb/mempool"
)

var (
	Bearer            string
	MevDynamoDBClient mempool_txs.MempoolTxDynamoDB
)

func InitMevDynamoDBClient(creds dynamodb_client.DynamoDBCredentials) {
	MevDynamoDBClient = mempool_txs.NewMempoolTxDynamoDB(creds)
}
