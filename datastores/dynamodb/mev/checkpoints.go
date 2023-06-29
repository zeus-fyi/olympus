package dynamodb_mev

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type CheckpointsDynamoDBTableKeys struct {
	CheckpointName string `dynamodbav:"checkpointName"`
}

var (
	MainnetCheckpointsDynamoDBTable = aws.String("MainnetEthereumCheckpoints")
)

type CheckpointsDynamoDB struct {
	CheckpointsDynamoDBTableKeys
	TTL int `dynamodbav:"ttl"`
}

func (m *MevDynamoDB) PutCheckpoint(ctx context.Context, checkpoint CheckpointsDynamoDB) error {
	now := time.Now()
	fourHours := now.Add(time.Hour * 4)
	unixTimestamp := fourHours.Unix()
	ttl := unixTimestamp
	checkpoint.TTL = int(ttl)
	item, err := attributevalue.MarshalMap(checkpoint)
	if err != nil {
		return err
	}
	_, err = m.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: MainnetCheckpointsDynamoDBTable,
	})
	if err != nil {
		return err
	}
	return err
}

func (m *MevDynamoDB) GetCheckpoint(ctx context.Context, checkpoint CheckpointsDynamoDB) (bool, error) {
	keymap, err := attributevalue.MarshalMap(checkpoint.CheckpointsDynamoDBTableKeys)
	if err != nil {
		return false, err
	}
	tableName := MainnetCheckpointsDynamoDBTable
	resp, err := m.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      tableName,
		Key:            keymap,
		ConsistentRead: aws.Bool(true),
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
