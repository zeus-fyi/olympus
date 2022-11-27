package auth_startup

import (
	"context"

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
	DoctlToken   string
	ArtemisEcdsaKeys
	Beacons
}

type ArtemisEcdsaKeys struct {
	Mainnet string
	Goerli  string
}

type Beacons struct {
	MainnetNodeUrl string
	GoerliNodeUrl  string
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

	log.Info().Msg("RunDigitalOceanS3BucketObjSecretsProcedure finished")

	doctlToken, err := s3SecretsReader.MemFS.ReadFile("secrets/doctl.txt")
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}

	sw := SecretsWrapper{}
	sw.DoctlToken = string(doctlToken)
	pgAuth, err := s3SecretsReader.MemFS.ReadFile("secrets/postgres-auth.txt")
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}
	sw.PostgresAuth = string(pgAuth)
	return s3SecretsReader.MemFS, sw
}

func RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Artemis: RunDigitalOceanS3BucketObjSecretsProcedure starting")

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
		log.Fatal().Msg("RunArtemisDigitalOceanS3BucketObjSecretsProcedure: MakeFile failed, shutting down the server")
		misc.DelayedPanic(err)
	}

	unzipDir := "./secrets"
	err = s3SecretsReader.DecryptAndUnGzipToInMemFs(&authCfg.Path, unzipDir)
	if err != nil {
		log.Fatal().Msg("RunArtemisDigitalOceanS3BucketObjSecretsProcedure: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}

	log.Info().Msg("RunArtemisDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	goerliNodeUrl, err := s3SecretsReader.MemFS.ReadFile("secrets/artemis-beacon-goerli.txt")
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}
	sw.GoerliNodeUrl = string(goerliNodeUrl)
	mainnetNodeUrl, err := s3SecretsReader.MemFS.ReadFile("secrets/artemis-beacon-mainnet.txt")
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}
	sw.MainnetNodeUrl = string(mainnetNodeUrl)
	artemisKey, err := s3SecretsReader.MemFS.ReadFile("secrets/artemis-eth-goerli.txt")
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}
	sw.ArtemisEcdsaKeys.Goerli = string(artemisKey)
	pgAuth, err := s3SecretsReader.MemFS.ReadFile("secrets/postgres-auth.txt")
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}
	sw.PostgresAuth = string(pgAuth)
	return s3SecretsReader.MemFS, sw
}
