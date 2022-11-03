package auth_startup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
)

type AuthStartupTestSuite struct {
	s3secrets.S3SecretsManagerTestSuite
}

// TestRead, you'll need to set the secret values to run the test
func (t *AuthStartupTestSuite) TestAuthStartup() {
	ctx := context.Background()

	keysCfg := AuthKeysCfg{
		AgePrivKey:    t.Tc.LocalAgePkey,
		AgePubKey:     t.Tc.LocalAgePubkey,
		SpacesKey:     t.Tc.LocalS3SpacesKey,
		SpacesPrivKey: t.Tc.LocalS3SpacesSecret,
	}
	authCfg := NewDefaultAuthClient(ctx, keysCfg)
	inMemFs := RunDigitalOceanS3BucketObjAuthProcedure(ctx, authCfg)

	t.Require().NotEmpty(inMemFs)
}

func TestAuthStartupTestSuite(t *testing.T) {
	suite.Run(t, new(AuthStartupTestSuite))
}
