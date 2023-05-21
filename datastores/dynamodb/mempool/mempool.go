package mempool_txs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rs/zerolog/log"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
)

type MempoolTxDynamoDB struct {
	*dynamodb.Client
}

func NewMempoolTxDynamoDB(creds dynamodb_client.DynamoDBCredentials) MempoolTxDynamoDB {
	d, err := dynamodb_client.NewDynamoDBClient(context.Background(), creds)
	if err != nil {
		log.Err(err)
	}
	return MempoolTxDynamoDB{
		d.Client,
	}
}

type MempoolTxDynamoDBTableKeys struct {
	Pubkey  string `dynamodbav:"pubkey"`
	TxOrder int    `dynamodbav:"txOrder"`
}

var (
	MainnetMempoolTxsTableName = aws.String("MempoolTxsMainnet")
	GoerliMempoolTxsTableName  = aws.String("MempoolTxsGoerli")
)

type MempoolTxsDynamoDB struct {
	MempoolTxDynamoDBTableKeys
	Tx  string `dynamodbav:"tx"`
	TTL int    `dynamodbav:"ttl"`
}

func (m *MempoolTxDynamoDB) GetMempoolTxs(ctx context.Context, network string) (*dynamodb.QueryOutput, error) {
	var mempoolTxsTableName *string
	if network == "mainnet" {
		mempoolTxsTableName = MainnetMempoolTxsTableName
	} else if network == "goerli" {
		mempoolTxsTableName = GoerliMempoolTxsTableName
	}
	r, err := m.Query(ctx, &dynamodb.QueryInput{
		TableName:                 mempoolTxsTableName,
		AttributesToGet:           nil,
		ConditionalOperator:       "",
		ConsistentRead:            nil,
		ExclusiveStartKey:         nil,
		ExpressionAttributeNames:  nil,
		ExpressionAttributeValues: nil,
		FilterExpression:          nil,
		IndexName:                 nil,
		KeyConditionExpression:    nil,
		KeyConditions:             nil,
		Limit:                     nil,
		ProjectionExpression:      nil,
		QueryFilter:               nil,
		ReturnConsumedCapacity:    "",
		ScanIndexForward:          nil,
		Select:                    "",
	})
	return r, err
}
