package artemis_hydra_orchestrations_auth

import (
	"context"
	aws_secrets "github.com/zeus-fyi/zeus/pkg/aegis/aws"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type ArtemisHydraSecretsManagerTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *ArtemisHydraSecretsManagerTestSuite) SetupTest() {
	s.InitLocalConfigs()
	auth := aws_secrets.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySecretManager,
		SecretKey: s.Tc.AwsSecretKeySecretManager,
	}
	InitHydraSecretManagerAuthAWS(ctx, auth)
}

func (s *ArtemisHydraSecretsManagerTestSuite) TestFetchServiceRoutesAuths() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	// TODO
	GetServiceRoutesAuths(ctx, ou)
}

func TestArtemisHydraSecretsManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisHydraSecretsManagerTestSuite))
}
