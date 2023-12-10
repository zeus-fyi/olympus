package aws_secrets

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

type AegisAwsSecretsTestSuite struct {
	s3secrets.S3SecretsManagerTestSuite
}

var ctx = context.Background()

// TestRead, you'll need to set the secret values to run the test
func (t *AegisAwsSecretsTestSuite) SetupTest() {
	t.InitLocalConfigs()
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: t.Tc.AwsAccessKeySecretManager,
		SecretKey: t.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, auth)
}

func (t *AegisAwsSecretsTestSuite) TestGetMockingbirdPlatformSecret() {
	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID

	svrl, err := GetMockingbirdPlatformSecrets(ctx, ou, "reddit")
	t.Require().Nil(err)
	t.Require().NotNil(svrl)

	svrl, err = GetMockingbirdPlatformSecrets(ctx, ou, "twitter")
	t.Require().Nil(err)
	t.Require().NotNil(svrl)

	svrl, err = GetMockingbirdPlatformSecrets(ctx, ou, "discord")
	t.Require().Nil(err)
	t.Require().NotNil(svrl)
}

func (t *AegisAwsSecretsTestSuite) TestSecretsRetrieval() {
	ou := org_users.OrgUser{}
	ou.OrgID = t.Tc.ProductionLocalTemporalOrgID
	ou.UserID = t.Tc.ProductionLocalTemporalUserID

	svrl, err := ReadSecretReferences(ctx, ou)
	t.Require().Nil(err)
	t.Require().NotNil(svrl)

	for _, svr := range svrl {
		fmt.Println(svr.Name, svr.Key)
	}
}

func TestAegisAwsSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(AegisAwsSecretsTestSuite))
}
