package auth_startup

import (
	"context"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

type SecretsWrapper struct {
	PostgresAuth string
}

func RunDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Zeus: RunDigitalOceanS3BucketObjSecretsProcedure starting")

	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String("secrets.tar.gz.age"),
	}
	authCfg.S3KeyValue = input
	s3Reader := s3reader.NewS3ClientReader(authCfg.s3BaseClient)
	s3SecretsReader := s3secrets.NewS3Secrets(authCfg.a, s3Reader)
	buf := s3SecretsReader.ReadBytes(ctx, &authCfg.Path, authCfg.S3KeyValue)

	tmpPath := filepaths.Path{}
	tmpPath.DirOut = "./"
	tmpPath.FnOut = "secrets.tar.gz.age"
	err := s3SecretsReader.MemFS.MakeFileIn(&authCfg.Path, buf.Bytes())
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: MakeFile failed, shutting down the server")
		misc.DelayedPanic(err)
	}

	unzipDir := "./secrets"
	err = s3SecretsReader.DecryptAndUnGzipToInMemFs(&authCfg.Path, unzipDir)
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}

	log.Info().Msg("Zeus: RunDigitalOceanS3BucketObjSecretsProcedure finished")

	doctlToken, err := s3SecretsReader.MemFS.ReadFile("secrets/doctl.txt")
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}
	cmd := exec.Command("doctl", "auth", "init", "-t", string(doctlToken))
	err = cmd.Run()
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: failed to auth doctl, shutting down the server")
		misc.DelayedPanic(err)
	}

	sw := SecretsWrapper{}
	pgAuth, err := s3SecretsReader.MemFS.ReadFile("secrets/postgres-auth.txt")
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}
	sw.PostgresAuth = string(pgAuth)
	return s3SecretsReader.MemFS, sw
}
