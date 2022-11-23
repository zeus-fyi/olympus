package auth_startup

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type AuthStartupTestSuite struct {
	s3secrets.S3SecretsManagerTestSuite
}

// TestRead, you'll need to set the secret values to run the test

func (t *AuthStartupTestSuite) TestSecretsEncrypt() {

	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./secrets",
		DirOut:      "./",
		FnIn:        "secrets",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}

	err := t.S3Secrets.GzipAndEncrypt(&p)
	t.Require().Nil(err)
}

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

	authCfg.Path.FnIn = "secrets.tar.gz.age"
	authCfg.Path.FnOut = "secrets.tar.gz"
	inMemSecrets, sw := RunDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
	t.Require().NotEmpty(inMemSecrets)
	t.Require().NotEmpty(sw)

	//b, err := inMemSecrets.ReadFile("secrets/doctl.txt")
	//t.Require().NotEmpty(b)
	//t.Require().Nil(err)
	//
	//token := string(b)
	//cmd := exec.Command("doctl", "auth", "init", "-t", token)
	//err = cmd.Run()

	t.Assert().Equal(t.Tc.ProdDbPgconn, sw.PostgresAuth)

}

func TestAuthStartupTestSuite(t *testing.T) {
	suite.Run(t, new(AuthStartupTestSuite))
}
