package dynamodb_web3signer

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/configs"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamoDB"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	"testing"
)

type Web3SignerDynamoDBTestSuite struct {
	Web3SignerDynamoDB
	test_suites_base.TestSuite
}

var (
	ctx    = context.Background()
	region = "us-west-1"
)

func (t *Web3SignerDynamoDBTestSuite) SetupTest() {
	t.Tc = configs.InitLocalTestConfigs()
	creds := dynamodb_client.DynamoDBCredentials{
		Region:       region,
		AccessKey:    t.Tc.AwsAccessKey,
		AccessSecret: t.Tc.AwsSecretKey,
	}
	d, err := NewWeb3SignerDynamoDB(ctx, creds)
	t.Require().Nil(err)
	t.Web3SignerDynamoDB = d
}

func (t *Web3SignerDynamoDBTestSuite) TestListTables() {
	tables, err := t.DynamoDB.ListTables(ctx, nil)
	t.Require().Nil(err)
	t.Require().NotNil(tables)
}

func TestWeb3SignerDynamoDBTestSuite(t *testing.T) {
	suite.Run(t, new(Web3SignerDynamoDBTestSuite))
}
