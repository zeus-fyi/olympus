package dynamodb_mev

import dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"

func (t *MevDynamoDBTestSuite) TestPutCheckpoint() {
	creds := dynamodb_client.DynamoDBCredentials{
		Region:       region,
		AccessKey:    t.Tc.AwsAccessKeyDynamoDB,
		AccessSecret: t.Tc.AwsSecretKeyDynamoDB,
	}
	m := NewMevDynamoDB(creds)

	ckp := CheckpointsDynamoDB{
		CheckpointsDynamoDBTableKeys: CheckpointsDynamoDBTableKeys{
			CheckpointName: "0x123",
		},
		Timestamp: 1,
	}
	err := m.PutCheckpoint(ctx, ckp)
	t.Require().Nil(err)
}

func (t *MevDynamoDBTestSuite) TestGetCheckpoint() {
	creds := dynamodb_client.DynamoDBCredentials{
		Region:       region,
		AccessKey:    t.Tc.AwsAccessKeyDynamoDB,
		AccessSecret: t.Tc.AwsSecretKeyDynamoDB,
	}
	m := NewMevDynamoDB(creds)
	ckp := CheckpointsDynamoDB{
		CheckpointsDynamoDBTableKeys: CheckpointsDynamoDBTableKeys{
			CheckpointName: "0x123",
		},
	}
	err := m.GetBlockCheckpointTime(ctx, &ckp)
	t.Require().Nil(err)
	t.Assert().NotEmpty(ckp.Timestamp)
}
