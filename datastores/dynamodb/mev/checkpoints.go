package dynamodb_mev

import (
	"context"

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
	Timestamp int `dynamodbav:"timestamp"`
	TTL       int `dynamodbav:"ttl"`
}

func (m *MevDynamoDB) PutCheckpoint(ctx context.Context, checkpoint CheckpointsDynamoDB) error {
	//now := time.Now()
	//fourHours := now.Add(time.Hour * 4)
	//unixTimestamp := fourHours.Unix()
	//ttl := unixTimestamp
	//checkpoint.TTL = int(ttl)
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

func (m *MevDynamoDB) GetBlockCheckpointTime(ctx context.Context, checkpoint *CheckpointsDynamoDB) error {
	keymap, err := attributevalue.MarshalMap(checkpoint.CheckpointsDynamoDBTableKeys)
	if err != nil {
		return err
	}
	tableName := MainnetCheckpointsDynamoDBTable
	resp, err := m.GetItem(ctx, &dynamodb.GetItemInput{
		Key:             keymap,
		TableName:       tableName,
		AttributesToGet: []string{"timestamp"},
		ConsistentRead:  aws.Bool(true),
	})
	if err != nil {
		return err
	}
	// If the item is found in the table, the response will contain the item.
	// If the item is not found, the response will be empty.
	err = attributevalue.UnmarshalMap(resp.Item, &checkpoint)
	if err != nil {
		return err
	}
	return err
}
