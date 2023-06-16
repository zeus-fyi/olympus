package dynamodb_client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rs/zerolog/log"
)

type DynamoDB struct {
	*dynamodb.Client
	DynamoDBCredentials
}

type DynamoDBCredentials struct {
	Region       string
	AccessKey    string
	AccessSecret string
}

func NewDynamoDBClient(ctx context.Context, creds DynamoDBCredentials) (DynamoDB, error) {
	d := DynamoDB{
		DynamoDBCredentials: creds,
	}
	err := d.InitDynamoDBClient(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err)
		return d, err
	}
	return d, err
}

func (d *DynamoDB) InitDynamoDBClient(ctx context.Context) error {
	creds := credentials.NewStaticCredentialsProvider(d.AccessKey, d.AccessSecret, "")
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(creds),
		config.WithRegion(d.Region),
	)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return err
	}
	dynDB := dynamodb.NewFromConfig(cfg)
	d.Client = dynDB
	return nil
}
