package mempool_txs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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

var (
	MainnetMempoolTxsTableName = aws.String("MempoolTxsMainnet")
	GoerliMempoolTxsTableName  = aws.String("MempoolTxsGoerli")
)

type MempoolTxDynamoDBTableKeys struct {
	Pubkey  string `dynamodbav:"pubkey"`
	TxOrder int    `dynamodbav:"txOrder"`
}

type MempoolTxsDynamoDB struct {
	MempoolTxDynamoDBTableKeys
	Tx  string `dynamodbav:"tx"`
	TTL int    `dynamodbav:"ttl"`
}

func (m *MempoolTxDynamoDB) GetMempoolTxs(ctx context.Context, network string) ([]MempoolTxsDynamoDB, error) {
	var mempoolTxsTableName *string
	if network == "mainnet" {
		mempoolTxsTableName = MainnetMempoolTxsTableName
	} else if network == "goerli" {
		mempoolTxsTableName = GoerliMempoolTxsTableName
	}
	fmt.Println(*mempoolTxsTableName)
	scanInput := &dynamodb.ScanInput{
		TableName: mempoolTxsTableName,
	}
	r, err := m.Scan(ctx, scanInput)
	if err != nil {
		return nil, err
	}

	var mempoolTxs []MempoolTxsDynamoDB
	err = attributevalue.UnmarshalListOfMaps(r.Items, &mempoolTxs)
	if err != nil {
		return nil, err
	}
	return mempoolTxs, nil
}
