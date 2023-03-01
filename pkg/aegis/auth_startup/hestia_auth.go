package auth_startup

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

const (
	secretsManagerAccessKey = "secrets/aws.secrets.manager.access.key.txt"
	secretsManagerSecretKey = "secrets/aws.secrets.manager.secret.key.txt"
)

func RunHestiaDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Hestia: RunDigitalOceanS3BucketObjSecretsProcedure starting")
	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("Hestia: RunDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.PostgresAuth = sw.ReadSecret(ctx, inMemSecrets, pgSecret)
	sw.BearerToken = sw.ReadSecret(ctx, inMemSecrets, temporalBearerSecret)
	sw.SecretsManagerAuthAWS.AccessKey = sw.ReadSecret(ctx, inMemSecrets, secretsManagerAccessKey)
	sw.SecretsManagerAuthAWS.SecretKey = sw.ReadSecret(ctx, inMemSecrets, secretsManagerSecretKey)
	log.Info().Msg("Hestia: RunDigitalOceanS3BucketObjSecretsProcedure succeeded")
	return inMemSecrets, sw
}
