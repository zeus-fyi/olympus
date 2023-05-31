package mempool_txs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rs/zerolog/log"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type MempoolTxDynamoDB struct {
	*dynamodb.Client
	*web3_client.Web3Client
}

func NewMempoolTxDynamoDB(creds dynamodb_client.DynamoDBCredentials, wc *web3_client.Web3Client) MempoolTxDynamoDB {
	d, err := dynamodb_client.NewDynamoDBClient(context.Background(), creds)
	if err != nil {
		log.Err(err)
	}
	return MempoolTxDynamoDB{
		d.Client, wc,
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
	scanInput := &dynamodb.ScanInput{
		TableName: mempoolTxsTableName,
	}
	r, err := m.Scan(ctx, scanInput)
	if err != nil {
		log.Err(err).Msg("GetMempoolTxs: error scanning mempool txs")
		return nil, err
	}

	var mempoolTxs []MempoolTxsDynamoDB
	err = attributevalue.UnmarshalListOfMaps(r.Items, &mempoolTxs)
	if err != nil {
		log.Err(err).Msg("GetMempoolTxs: error UnmarshalListOfMaps mempool txs")
		return nil, err
	}

	//fmt.Println("startingTxCount", len(mempoolTxs))
	//fmt.Println("endFilteredTxCount", len(txMap))
	//fmt.Println("filteredCount", len(mempoolTxs)-len(txMap))
	return mempoolTxs, nil
}

func (m *MempoolTxDynamoDB) RemoveMempoolTx(ctx context.Context, tx MempoolTxsDynamoDB) error {
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
