package dynamodb_web3signer

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rs/zerolog/log"
)

var AttestationsTableName = aws.String("Attestations")

type AttestationsDynamoDB struct {
	Web3SignerDynamoDBTableKeys
	SourceEpoch     int    `dynamodbav:"sourceEpoch"`
	TargetEpoch     int    `dynamodbav:"targetEpoch"`
	AttestationData string `dynamodbav:"attestationData"`
	TTL             int    `dynamodbav:"ttl"`
}

func (w *Web3SignerDynamoDB) PutAttestation(ctx context.Context, att AttestationsDynamoDB) error {
	now := time.Now()
	oneHourLater := now.Add(time.Hour)
	unixTimestamp := oneHourLater.Unix()
	ttl := unixTimestamp
	att.TTL = int(ttl)
	item, err := attributevalue.MarshalMap(att)
	if err != nil {
		log.Ctx(ctx).Error().Interface("att", att).Err(err).Msg("failed to marshal attestation")
		return err
	}
	_, err = w.DynamoDB.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: AttestationsTableName,
	})
	if err != nil {
		log.Ctx(ctx).Error().Interface("item", item).Err(err).Msg("failed to put attestation")
		return err
	}
	return err
}

func (w *Web3SignerDynamoDB) GetAttestation(ctx context.Context, tableKeys Web3SignerDynamoDBTableKeys) (AttestationsDynamoDB, error) {
	att := AttestationsDynamoDB{}
	item, err := attributevalue.MarshalMap(tableKeys)
	if err != nil {
		log.Ctx(ctx).Error().Interface("tableKeys", tableKeys).Err(err).Msg("failed to marshal tableKeys")
		return att, err
	}
	resp, err := w.DynamoDB.GetItem(ctx, &dynamodb.GetItemInput{
		Key:            item,
		TableName:      AttestationsTableName,
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		log.Ctx(ctx).Error().Interface("resp", resp).Err(err).Msg("failed to get attestation")
		return att, err
	}
	err = attributevalue.UnmarshalMap(resp.Item, &att)
	if err != nil {
		log.Ctx(ctx).Error().Interface("respItem", resp.Item).Err(err).Msg("failed to unmarshall last attestation")
		return att, err
	}
	return att, err
}
