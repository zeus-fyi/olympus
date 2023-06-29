package dynamodb_mev

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/configs"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type MevDynamoDBTestSuite struct {
	m MempoolTxsDynamoDB
	test_suites_base.TestSuite
}

var (
	ctx    = context.Background()
	region = "us-west-1"
)

func (t *MevDynamoDBTestSuite) SetupTest() {
	t.Tc = configs.InitLocalTestConfigs()
}

func (t *MevDynamoDBTestSuite) TestGetMempoolTxs() {
	creds := dynamodb_client.DynamoDBCredentials{
		Region:       region,
		AccessKey:    t.Tc.AwsAccessKeyDynamoDB,
		AccessSecret: t.Tc.AwsSecretKeyDynamoDB,
	}
	m := NewMevDynamoDB(creds)
	memTxs, err := m.GetMempoolTxs(ctx, "mainnet")
	t.Require().Nil(err)
	t.Require().NotNil(memTxs)
}

func TestMevDynamoDBTestSuite(t *testing.T) {
	suite.Run(t, new(MevDynamoDBTestSuite))
}
