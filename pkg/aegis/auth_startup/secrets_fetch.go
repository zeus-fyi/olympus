package auth_startup

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

const pgSecret = "secrets/postgres-auth.txt"
const doctlSecret = "secrets/doctl.txt"
const encryptedSecret = "secrets.tar.gz.age"
const secretBucketName = "zeus-fyi"

type SecretsWrapper struct {
	PostgresAuth string
	DoctlToken   string
	TemporalAuth temporal_auth.TemporalAuth
}

var secretsBucket = &s3.GetObjectInput{
	Bucket: aws.String(secretBucketName),
	Key:    aws.String(encryptedSecret),
}

func (s *SecretsWrapper) ReadSecret(ctx context.Context, inMemSecrets memfs.MemFS, fileName string) string {
	secret, err := inMemSecrets.ReadFile(fileName)
	if err != nil {
		log.Ctx(ctx).Fatal().Msgf("SecretsWrapper: ReadSecret failed, shutting down the server: %s", fileName)
		misc.DelayedPanic(err)
	}
	return string(secret)
}

func ReadEncryptedSecretsData(ctx context.Context, authCfg AuthConfig) memfs.MemFS {
	authCfg.S3KeyValue = secretsBucket
	s3Reader := s3reader.NewS3ClientReader(authCfg.s3BaseClient)
	s3SecretsReader := s3secrets.NewS3Secrets(authCfg.a, s3Reader)
	buf := s3SecretsReader.ReadBytes(ctx, &authCfg.Path, authCfg.S3KeyValue)

	tmpPath := filepaths.Path{}
	tmpPath.DirOut = "./"
	tmpPath.FnOut = encryptedSecret
	err := s3SecretsReader.MemFS.MakeFileIn(&authCfg.Path, buf.Bytes())
	if err != nil {
		log.Ctx(ctx).Fatal().Msg("ReadEncryptedSecretsData: MakeFile failed, shutting down the server")
		misc.DelayedPanic(err)
	}
	unzipDir := "./secrets"
	err = s3SecretsReader.DecryptAndUnGzipToInMemFs(&authCfg.Path, unzipDir)
	if err != nil {
		log.Ctx(ctx).Fatal().Msg("ReadEncryptedSecretsData: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}
	return s3SecretsReader.MemFS
}

func RunDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Zeus: RunDigitalOceanS3BucketObjSecretsProcedure starting")

	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("RunDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.DoctlToken = sw.ReadSecret(ctx, inMemSecrets, doctlSecret)
	sw.PostgresAuth = sw.ReadSecret(ctx, inMemSecrets, pgSecret)
	return inMemSecrets, sw
}

func RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Artemis: RunDigitalOceanS3BucketObjSecretsProcedure starting")
	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("RunArtemisDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.PostgresAuth = sw.ReadSecret(ctx, inMemSecrets, pgSecret)
	log.Info().Msg("RunArtemisDigitalOceanS3BucketObjSecretsProcedure succeeded")
	return inMemSecrets, sw
}
