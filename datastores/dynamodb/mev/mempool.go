package dynamodb_mev

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rs/zerolog/log"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
)

type MevDynamoDB struct {
	*dynamodb.Client
}

func NewMevDynamoDB(creds dynamodb_client.DynamoDBCredentials) MevDynamoDB {
	d, err := dynamodb_client.NewDynamoDBClient(context.Background(), creds)
	if err != nil {
		log.Err(err)
	}
	return MevDynamoDB{
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

func (m *MevDynamoDB) GetMempoolTxs(ctx context.Context, network string) ([]MempoolTxsDynamoDB, error) {
	var mempoolTxsTableName *string
	if network == "mainnet" {
		mempoolTxsTableName = MainnetMempoolTxsTableName
	} else if network == "goerli" {
		mempoolTxsTableName = GoerliMempoolTxsTableName
	}
	scanInput := &dynamodb.ScanInput{
		TableName: mempoolTxsTableName,
	}
	r, err := m.Scan(ctx, scanInput)
	if err != nil {
		log.Err(err).Msg("GetDynamoDBMempoolTxs: error scanning mempool txs")
		return nil, err
	}

	var mempoolTxs []MempoolTxsDynamoDB
	err = attributevalue.UnmarshalListOfMaps(r.Items, &mempoolTxs)
	if err != nil {
		log.Err(err).Msg("GetDynamoDBMempoolTxs: error UnmarshalListOfMaps mempool txs")
		return nil, err
	}

	//fmt.Println("startingTxCount", len(mempoolTxs))
	//fmt.Println("endFilteredTxCount", len(txMap))
	//fmt.Println("filteredCount", len(mempoolTxs)-len(txMap))
	return mempoolTxs, nil
}

func (m *MevDynamoDB) RemoveMempoolTx(ctx context.Context, tx MempoolTxsDynamoDB) error {
	keymap, err := attributevalue.MarshalMap(tx.MempoolTxDynamoDBTableKeys)
	if err != nil {
		return err
	}
	_, err = m.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: MainnetMempoolTxsTableName,
		Key:       keymap,
	})
	if err != nil {
		log.Err(err).Msg("RemoveMempoolTx: error deleting mempool tx")
		return nil
	}
	return nil
}
