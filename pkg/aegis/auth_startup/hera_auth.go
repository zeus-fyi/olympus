package auth_startup

import (
	"context"

	"github.com/rs/zerolog/log"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

const (
	heraOpenAIAuth = "secrets/hera.openai.auth.txt"
)

func RunHeraDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Hera: RunDigitalOceanS3BucketObjSecretsProcedure starting")
	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("Hera: RunHeraDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := InitHera(ctx, inMemSecrets)
	log.Info().Msg("Hera: RunHeraDigitalOceanS3BucketObjSecretsProcedure succeeded")
	return inMemSecrets, sw
}

func InitHera(ctx context.Context, inMemSecrets memfs.MemFS) SecretsWrapper {
	log.Info().Msg("Hera: InitHera starting")
	secrets := SecretsWrapper{}
	secrets.PostgresAuth = secrets.MustReadSecret(ctx, inMemSecrets, PgSecret)
	secrets.OpenAIToken = secrets.MustReadSecret(ctx, inMemSecrets, heraOpenAIAuth)
	hera_openai.InitHeraOpenAI(secrets.OpenAIToken)
	log.Info().Msg("Hera: InitHera done")
	return secrets
}
