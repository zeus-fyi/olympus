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
	sw.PostgresAuth = sw.ReadSecret(ctx, inMemSecrets, pgSecret)
	sw.AccessKeyHydraDynamoDB = sw.ReadSecret(ctx, inMemSecrets, hydraAccessKeyDynamoDB)
	sw.SecretKeyHydraDynamoDB = sw.ReadSecret(ctx, inMemSecrets, hydraSecretKeyDynamoDB)
	log.Info().Msg("Hydra: RunHydraDigitalOceanS3BucketObjSecretsProcedure succeeded")
	return inMemSecrets, sw
}
