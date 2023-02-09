package dynamodb_client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type DynamoDBTestSuite struct {
	DynamoDB
	test_suites_base.TestSuite
}

var (
	ctx    = context.Background()
	region = "us-west-1"
)

func (t *DynamoDBTestSuite) SetupTest() {
	t.Tc = configs.InitLocalTestConfigs()
	d, err := NewDynamoDBClient(ctx, region, t.Tc.AwsAccessKey, t.Tc.AwsSecretKey)
	t.Require().Nil(err)
	t.DynamoDB = d
}

func (t *DynamoDBTestSuite) TestListTables() {
	tables, err := t.DynamoDB.ListTables(ctx, nil)
	t.Require().Nil(err)
	t.Require().NotNil(tables)
}

func TestDynamoDBTestSuite(t *testing.T) {
	suite.Run(t, new(DynamoDBTestSuite))
}
