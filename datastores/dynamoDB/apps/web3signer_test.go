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
	ctx               = context.Background()
	region            = "us-west-1"
	dynamoDBTableKeys = Web3SignerDynamoDBTableKeys{
		Pubkey:  "0x8a7addbf2857a72736205d861169c643545283a74a1ccb71c95dd2c9652acb89de226ca26d60248c4ef9591d7e010288",
		Network: "ephemery",
	}
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

func (t *Web3SignerDynamoDBTestSuite) TestAttestations() {
	att := AttestationsDynamoDB{
		Web3SignerDynamoDBTableKeys: dynamoDBTableKeys,
		SourceEpoch:                 10,
		TargetEpoch:                 11,
	}

	err := t.Web3SignerDynamoDB.PutAttestation(ctx, att)
	t.Require().Nil(err)

	returnedAtt, err := t.Web3SignerDynamoDB.GetAttestation(ctx, dynamoDBTableKeys)
	t.Require().Nil(err)
	t.Require().NotEmpty(returnedAtt)
	t.Require().Equal(att, returnedAtt)
}

func (t *Web3SignerDynamoDBTestSuite) TestBlockProposals() {
	bp := BlockProposalsDynamoDB{
		Web3SignerDynamoDBTableKeys: dynamoDBTableKeys,
		Slot:                        10,
	}
	err := t.Web3SignerDynamoDB.PutBlockProposal(ctx, bp)
	t.Require().Nil(err)

	returnedBP, err := t.Web3SignerDynamoDB.GetBlockProposal(ctx, dynamoDBTableKeys)
	t.Require().Nil(err)
	t.Require().NotNil(returnedBP)
	t.Require().Equal(bp, returnedBP)
}

func (t *Web3SignerDynamoDBTestSuite) TestListTables() {
	tables, err := t.DynamoDB.ListTables(ctx, nil)
	t.Require().Nil(err)
	t.Require().NotNil(tables)
}

func TestWeb3SignerDynamoDBTestSuite(t *testing.T) {
	suite.Run(t, new(Web3SignerDynamoDBTestSuite))
}
