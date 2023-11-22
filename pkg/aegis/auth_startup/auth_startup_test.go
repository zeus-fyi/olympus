package auth_startup

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/suite"
	s3uploader "github.com/zeus-fyi/olympus/datastores/s3/upload"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type AuthStartupTestSuite struct {
	s3secrets.S3SecretsManagerTestSuite
}

var ctx = context.Background()

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

	fn := "secrets.tar.gz.age"
	bucketName := "zeus-fyi"
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fn),
	}

	uploader := s3uploader.NewS3ClientUploader(t.S3)
	err = uploader.Upload(ctx, p, input)
	t.Require().Nil(err)
}

func (t *AuthStartupTestSuite) TestAuthStartup() {
	keysCfg := auth_keys_config.AuthKeysCfg{
		AgePrivKey:    t.Tc.LocalAgePkey,
		AgePubKey:     t.Tc.LocalAgePubkey,
		SpacesKey:     t.Tc.LocalS3SpacesKey,
		SpacesPrivKey: t.Tc.LocalS3SpacesSecret,
	}
	authCfg := NewDefaultAuthClient(ctx, keysCfg)
	inMemFs := ReadEncryptedSecretsData(ctx, authCfg)

	t.Require().NotEmpty(inMemFs)
	ra := RedditAuthConfig{}
	sw := SecretsWrapper{}

	sb := sw.ReadSecretBytes(ctx, inMemFs, redditSecretsJson)
	err := json.Unmarshal(sb, &ra)
	if err != nil {
		panic(err)
	}
	//authCfg.Path.FnIn = "secrets.tar.gz.age"
	//authCfg.Path.FnOut = "secrets.tar.gz"
	//inMemSecrets, sw := RunDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
	//t.Require().NotEmpty(inMemSecrets)
	//t.Require().NotEmpty(sw)

	//b, err := inMemSecrets.ReadFile("secrets/doctl.txt")
	//t.Require().NotEmpty(b)
	//t.Require().Nil(err)
	//
	//token := string(b)
	//cmd := exec.Command("doctl", "auth", "init", "-t", token)
	//err = cmd.Run()

	//t.Assert().Equal(t.Tc.ProdDbPgconn, sw.PostgresAuth)

}
func (t *AuthStartupTestSuite) TestHestiaAuthStartup() {
	keysCfg := auth_keys_config.AuthKeysCfg{
		AgePrivKey:    t.Tc.LocalAgePkey,
		AgePubKey:     t.Tc.LocalAgePubkey,
		SpacesKey:     t.Tc.LocalS3SpacesKey,
		SpacesPrivKey: t.Tc.LocalS3SpacesSecret,
	}
	authCfg := NewDefaultAuthClient(ctx, keysCfg)
	inMemFs := ReadEncryptedSecretsData(ctx, authCfg)

	t.Require().NotEmpty(inMemFs)

	fs, sw := RunHestiaDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
	b, err := fs.ReadFile("secrets/artemis.ethereum.goerli.beacon.txt")
	t.Require().NotEmpty(b)
	t.Require().Nil(err)

	InitArtemisEthereum(ctx, fs, sw)
	//token := string(b)
	//cmd := exec.Command("doctl", "auth", "init", "-t", token)
	//err = cmd.Run()

	t.Assert().Equal(t.Tc.ProdDbPgconn, sw.PostgresAuth)
}
func (t *AuthStartupTestSuite) TestArtemisAuthStartup() {
	keysCfg := auth_keys_config.AuthKeysCfg{
		AgePrivKey:    t.Tc.LocalAgePkey,
		AgePubKey:     t.Tc.LocalAgePubkey,
		SpacesKey:     t.Tc.LocalS3SpacesKey,
		SpacesPrivKey: t.Tc.LocalS3SpacesSecret,
	}
	authCfg := NewDefaultAuthClient(ctx, keysCfg)
	inMemFs := ReadEncryptedSecretsData(ctx, authCfg)

	t.Require().NotEmpty(inMemFs)

	fs, sw := RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
	b, err := fs.ReadFile("secrets/artemis.ethereum.goerli.beacon.txt")
	t.Require().NotEmpty(b)
	t.Require().Nil(err)

	InitArtemisEthereum(ctx, fs, sw)
	//token := string(b)
	//cmd := exec.Command("doctl", "auth", "init", "-t", token)
	//err = cmd.Run()

	t.Assert().Equal(t.Tc.ProdDbPgconn, sw.PostgresAuth)
}

func TestAuthStartupTestSuite(t *testing.T) {
	suite.Run(t, new(AuthStartupTestSuite))
}
