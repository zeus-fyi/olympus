package artemis_orchestration_auth

import (
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	dynamodb_mev "github.com/zeus-fyi/olympus/datastores/dynamodb/mev"
)

var (
	Bearer            string
	MevDynamoDBClient dynamodb_mev.MevDynamoDB
)

func InitMevDynamoDBClient(creds dynamodb_client.DynamoDBCredentials) {
	MevDynamoDBClient = dynamodb_mev.NewMevDynamoDB(creds)
}
