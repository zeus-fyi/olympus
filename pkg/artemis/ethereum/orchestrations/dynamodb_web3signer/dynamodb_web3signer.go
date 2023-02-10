package dynamodb_web3signer_client

import (
	"context"
	"github.com/rs/zerolog/log"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	dynamodb_web3signer "github.com/zeus-fyi/olympus/datastores/dynamodb/apps"
)

var Web3SignerDynamoDBClient dynamodb_web3signer.Web3SignerDynamoDB

func InitWeb3SignerDynamoDBClient(ctx context.Context, creds dynamodb_client.DynamoDBCredentials) {
	var err error
	Web3SignerDynamoDBClient, err = dynamodb_web3signer.NewWeb3SignerDynamoDB(ctx, creds)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("failed to init Web3SignerDynamoDBClient")
		panic(err)
	}
}
