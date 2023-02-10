package dynamodb_web3signer

import (
	"context"
	"github.com/rs/zerolog/log"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
)

type Web3SignerDynamoDB struct {
	dynamodb_client.DynamoDB
}

type Web3SignerDynamoDBTableKeys struct {
	Pubkey  string `dynamodbav:"pubkey"`
	Network string `dynamodbav:"network"`
}

func NewWeb3SignerDynamoDB(ctx context.Context, creds dynamodb_client.DynamoDBCredentials) (Web3SignerDynamoDB, error) {
	w3DynDB := Web3SignerDynamoDB{}
	client, err := dynamodb_client.NewDynamoDBClient(ctx, creds)
	if err != nil {
		log.Ctx(ctx).Error().Err(err)
		return w3DynDB, err
	}
	w3DynDB.DynamoDB = client
	return w3DynDB, nil
}
