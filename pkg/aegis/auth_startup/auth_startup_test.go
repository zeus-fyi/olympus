package auth_startup

import (
	"context"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
)

type AuthStartupTestSuite struct {
	s3secrets.S3SecretsManagerTestSuite
}

// TestRead, you'll need to set the secret values to run the test
func (t *AuthStartupTestSuite) TestAuthStartup() {
	ctx := context.Background()

	keysCfg := auth_keys_config.AuthKeysCfg{
		AgePrivKey:    t.Tc.LocalAgePkey,
		AgePubKey:     t.Tc.LocalAgePubkey,
		SpacesKey:     t.Tc.LocalS3SpacesKey,
		SpacesPrivKey: t.Tc.LocalS3SpacesSecret,
	}
	authCfg := NewDefaultAuthClient(ctx, keysCfg)
	inMemFs := RunDigitalOceanS3BucketObjAuthProcedure(ctx, authCfg)

	t.Require().NotEmpty(inMemFs)

	authCfg.Path.Fn = "secrets.tar.gz.age"
	authCfg.Path.FnOut = "secrets.tar.gz"
	inMemSecrets := RunDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
	t.Require().NotEmpty(inMemSecrets)

	b, err := inMemSecrets.ReadFile("secrets/doctl.txt")
	t.Require().NotEmpty(b)
	t.Require().Nil(err)

	token := string(b)
	cmd := exec.Command("doctl", "auth", "init", "-t", token)
	err = cmd.Run()

	t.Require().Nil(err)

}

func TestAuthStartupTestSuite(t *testing.T) {
	suite.Run(t, new(AuthStartupTestSuite))
}
