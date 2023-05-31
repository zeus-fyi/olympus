package auth_startup

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

const (
	hydraAccessKeyDynamoDB = "secrets/hydra.dynamodb.access.key.txt"
	hydraSecretKeyDynamoDB = "secrets/hydra.dynamodb.secret.key.txt"
)

func RunHydraDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Hydra: RunHydraDigitalOceanS3BucketObjSecretsProcedure starting")
	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("Hydra: RunHydraDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.PostgresAuth = sw.MustReadSecret(ctx, inMemSecrets, pgSecret)
	sw.AccessKeyHydraDynamoDB = sw.MustReadSecret(ctx, inMemSecrets, hydraAccessKeyDynamoDB)
	sw.SecretKeyHydraDynamoDB = sw.MustReadSecret(ctx, inMemSecrets, hydraSecretKeyDynamoDB)
	sw.PagerDutyApiKey = sw.MustReadSecret(ctx, inMemSecrets, pagerDutySecret)
	sw.PagerDutyRoutingKey = sw.MustReadSecret(ctx, inMemSecrets, pagerDutyRoutingKey)

	sw.SecretsManagerAuthAWS.AccessKey = sw.MustReadSecret(ctx, inMemSecrets, secretsManagerAccessKey)
	sw.SecretsManagerAuthAWS.SecretKey = sw.MustReadSecret(ctx, inMemSecrets, secretsManagerSecretKey)
	log.Info().Msg("Hydra: RunHydraDigitalOceanS3BucketObjSecretsProcedure succeeded")
	return inMemSecrets, sw
}
