package dynamodb_mev

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type TxBlacklistDynamoDBTableKeys struct {
	TxHash string `dynamodbav:"txHash"`
}

var (
	MainnetTxBlacklistTableName = aws.String("MempoolTxsBlacklistMainnet")
)

type TxBlacklistDynamoDB struct {
	TxBlacklistDynamoDBTableKeys
	TTL int `dynamodbav:"ttl"`
}

func (m *MevDynamoDB) PutTxBlacklist(ctx context.Context, txBlacklist TxBlacklistDynamoDB) error {
	now := time.Now()
	fourHours := now.Add(time.Hour * 4)
	unixTimestamp := fourHours.Unix()
	ttl := unixTimestamp
	txBlacklist.TTL = int(ttl)
	item, err := attributevalue.MarshalMap(txBlacklist)
	if err != nil {
		return err
	}
	_, err = m.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: MainnetTxBlacklistTableName,
	})
	if err != nil {
		return err
	}
	return err
}

func (m *MevDynamoDB) GetTxBlacklist(ctx context.Context, txBlacklist TxBlacklistDynamoDB) (bool, error) {
	keymap, err := attributevalue.MarshalMap(txBlacklist.TxBlacklistDynamoDBTableKeys)
	if err != nil {
		return false, err
	}
	tableName := MainnetTxBlacklistTableName
	resp, err := m.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      tableName,
		Key:            keymap,
		ConsistentRead: aws.Bool(false),
	})
	if err != nil {
		return false, err
	}
	// If the item is found in the table, the response will contain the item.
	// If the item is not found, the response will be empty.
	if len(resp.Item) > 0 {
		return true, nil
	}
	return false, err
}
