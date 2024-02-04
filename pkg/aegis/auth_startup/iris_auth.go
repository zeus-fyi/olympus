package auth_startup

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

func RunIrisDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Iris: RunIrisDigitalOceanS3BucketObjSecretsProcedure starting")
	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("Iris: RunIrisDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.PostgresAuth = sw.MustReadSecret(ctx, inMemSecrets, PgSecret)
	sw.BearerToken = sw.MustReadSecret(ctx, inMemSecrets, temporalBearerSecret)
	sw.AccessKeyHydraDynamoDB = sw.MustReadSecret(ctx, inMemSecrets, HydraAccessKeyDynamoDB)
	sw.SecretKeyHydraDynamoDB = sw.MustReadSecret(ctx, inMemSecrets, HydraSecretKeyDynamoDB)
	sw.ZeroXApiKey = sw.MustReadSecret(ctx, inMemSecrets, zeroXApiKey)
	sw.SendGridAPIKey = sw.MustReadSecret(ctx, inMemSecrets, sendGridAPIKey)

	sw.SecretsManagerAuthAWS.AccessKey = sw.MustReadSecret(ctx, inMemSecrets, secretsManagerAccessKey)
	sw.SecretsManagerAuthAWS.SecretKey = sw.MustReadSecret(ctx, inMemSecrets, secretsManagerSecretKey)
	log.Info().Msg("Iris: RunIrisDigitalOceanS3BucketObjSecretsProcedure succeeded")
	return inMemSecrets, sw
}
