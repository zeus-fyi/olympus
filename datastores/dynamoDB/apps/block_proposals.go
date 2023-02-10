package dynamodb_web3signer

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var BlockProposalsTableName = aws.String("BlockProposals")

type BlockProposalsDynamoDB struct {
	Web3SignerDynamoDBTableKeys
	Slot int `dynamodbav:"slot"`
}

func (w *Web3SignerDynamoDB) PutBlockProposal(ctx context.Context, bp BlockProposalsDynamoDB) error {
	item, err := attributevalue.MarshalMap(bp)
	if err != nil {
		log.Ctx(ctx).Error().Interface("att", bp).Err(err).Msg("failed to marshal block proposal")
		return err
	}
	_, err = w.DynamoDB.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: BlockProposalsTableName,
	})
	if err != nil {
		log.Ctx(ctx).Error().Interface("item", item).Err(err).Msg("failed to put block proposal")
		return err
	}
	return err
}

func (w *Web3SignerDynamoDB) GetBlockProposal(ctx context.Context, tableKeys Web3SignerDynamoDBTableKeys) (BlockProposalsDynamoDB, error) {
	bp := BlockProposalsDynamoDB{}
	keyMap, err := attributevalue.MarshalMap(tableKeys)
	if err != nil {
		log.Ctx(ctx).Error().Interface("keyMap", tableKeys).Err(err).Msg("failed to marshal tableKeys")
		return bp, err
	}
	resp, err := w.DynamoDB.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      BlockProposalsTableName,
		Key:            keyMap,
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		log.Ctx(ctx).Error().Interface("resp", resp).Err(err).Msg("failed to get last block proposal")
		return bp, err
	}
	err = attributevalue.UnmarshalMap(resp.Item, &bp)
	if err != nil {
		log.Ctx(ctx).Error().Interface("respItem", resp.Item).Err(err).Msg("failed to unmarshall last block proposal")
		return bp, err
	}
	return bp, err
}
