package auth_startup

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

const (
	apolloMainnetBeacon = "secrets/apollo.ethereum.mainnet.beacon.txt"
	apolloPostgresAuth  = "secrets/apollo.ethereum.postgres.auth.txt"

	apolloAtlassianOrgId  = "secrets/apollo.atlassian.org.id.txt"
	apolloAtlassianApiKey = "secrets/apollo.atlassian.api.key.txt"
)

func RunApolloDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Artemis: RunDigitalOceanS3BucketObjSecretsProcedure starting")
	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("RunArtemisDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := InitApolloEthereum(ctx, inMemSecrets)
	log.Info().Msg("RunArtemisDigitalOceanS3BucketObjSecretsProcedure succeeded")
	return inMemSecrets, sw
}

func InitApolloEthereum(ctx context.Context, inMemSecrets memfs.MemFS) SecretsWrapper {
	log.Info().Msg("Apollo: InitApolloEthereum starting")
	secrets := SecretsWrapper{}
	secrets.MainnetBeaconURL = secrets.MustReadSecret(ctx, inMemSecrets, apolloMainnetBeacon)
	secrets.PostgresAuth = secrets.MustReadSecret(ctx, inMemSecrets, apolloPostgresAuth)
	secrets.AegisPostgresAuth = secrets.MustReadSecret(ctx, inMemSecrets, PgSecret)
	log.Info().Msg("Apollo: InitApolloEthereum done")
	return secrets
}

func InitAtlassianKeys(ctx context.Context, inMemSecrets memfs.MemFS, secrets *SecretsWrapper) {
	if secrets == nil {
		secrets = &SecretsWrapper{}
	}
	secrets.AtlassianOrgId = secrets.MustReadSecret(ctx, inMemSecrets, apolloAtlassianOrgId)
	secrets.AtlassianApiKey = secrets.MustReadSecret(ctx, inMemSecrets, apolloAtlassianApiKey)
}
