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
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

const (
	temporalBearerSecret = "secrets/temporal.bearer.txt"
	pgSecret             = "secrets/postgres-auth.txt"
	doctlSecret          = "secrets/doctl.txt"
	rcloneSecret         = "secrets/rclone.conf"
	encryptedSecret      = "secrets.tar.gz.age"
	secretBucketName     = "zeus-fyi"
	pagerDutySecret      = "secrets/pagerduty.txt"
	pagerDutyRoutingKey  = "secrets/pagerduty.routing.key.txt"
)

type SecretsWrapper struct {
	PostgresAuth           string
	AegisPostgresAuth      string
	DoctlToken             string
	MainnetBeaconURL       string
	BearerToken            string
	OpenAIToken            string
	AccessKeyHydraDynamoDB string
	SecretKeyHydraDynamoDB string
	PagerDutyApiKey        string
	PagerDutyRoutingKey    string

	SecretsManagerAuthAWS aegis_aws_auth.AuthAWS
	TemporalAuth          temporal_auth.TemporalAuth
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
	log.Info().Msg("Artemis: RunArtemisDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.PostgresAuth = sw.ReadSecret(ctx, inMemSecrets, pgSecret)
	log.Info().Msg("Artemis: RunArtemisDigitalOceanS3BucketObjSecretsProcedure succeeded")
	return inMemSecrets, sw
}

func RunPoseidonDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Poseidon: RunDigitalOceanS3BucketObjSecretsProcedure starting")
	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("RunPoseidonDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.PostgresAuth = sw.ReadSecret(ctx, inMemSecrets, pgSecret)
	sw.BearerToken = sw.ReadSecret(ctx, inMemSecrets, temporalBearerSecret)
	log.Info().Msg("RunPoseidonDigitalOceanS3BucketObjSecretsProcedure succeeded")
	return inMemSecrets, sw
}

func RunAthenaDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Athena: RunDigitalOceanS3BucketObjSecretsProcedure starting")

	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("Athena: RunDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.PostgresAuth = sw.ReadSecret(ctx, inMemSecrets, pgSecret)

	p := filepaths.Path{
		PackageName: "",
		DirIn:       "",
		DirOut:      "/root/.config/rclone",
		FnOut:       "rclone.conf",
		Env:         "",
		Metadata:    nil,
		FilterFiles: string_utils.FilterOpts{},
	}
	rcloneConf, err := inMemSecrets.ReadFile(rcloneSecret)
	if err != nil {
		log.Err(err).Msg("Athena:  RunAthenaDigitalOceanS3BucketObjSecretsProcedure failed to set rclone conf")
		misc.DelayedPanic(err)
	}
	err = p.WriteToFileOutPath(rcloneConf)
	if err != nil {
		log.Err(err).Msg("Athena:  RunAthenaDigitalOceanS3BucketObjSecretsProcedure failed to set rclone conf")
		misc.DelayedPanic(err)
	}
	return inMemSecrets, sw
}
