package mempool_txs

import (
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
)

func (t *MempoolTxDynamoDBTestSuite) TestPutBlacklistTx() {
	creds := dynamodb_client.DynamoDBCredentials{
		Region:       region,
		AccessKey:    t.Tc.AwsAccessKeyDynamoDB,
		AccessSecret: t.Tc.AwsSecretKeyDynamoDB,
	}
	m := NewMempoolTxDynamoDB(creds, nil)

	txBlacklist := TxBlacklistDynamoDB{
		TxBlacklistDynamoDBTableKeys: TxBlacklistDynamoDBTableKeys{
			TxHash: "0x123",
		},
	}
	err := m.PutTxBlacklist(ctx, txBlacklist)
	t.Require().Nil(err)
}

func (t *MempoolTxDynamoDBTestSuite) TestGetBlacklistTx() {
	creds := dynamodb_client.DynamoDBCredentials{
		Region:       region,
		AccessKey:    t.Tc.AwsAccessKeyDynamoDB,
		AccessSecret: t.Tc.AwsSecretKeyDynamoDB,
	}
	m := NewMempoolTxDynamoDB(creds, nil)
	txBlacklist := TxBlacklistDynamoDB{
		TxBlacklistDynamoDBTableKeys: TxBlacklistDynamoDBTableKeys{
			TxHash: "0x123",
		},
	}
	found, err := m.GetTxBlacklist(ctx, txBlacklist)
	t.Require().Nil(err)
	t.Require().True(found)
}
