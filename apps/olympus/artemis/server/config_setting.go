package artemis_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_ethereum_transcations "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/transcations"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
)

var temporalProdAuthConfig = temporal_auth.TemporalAuth{
	ClientCertPath:   "/etc/ssl/certs/ca.pem",
	ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
	Namespace:        "production-artemis.ngb72",
	HostPort:         "production-artemis.ngb72.tmprl.cloud:7233",
}

func SetConfigByEnv(ctx context.Context, env string) {

	switch env {
	case "production":
		log.Info().Msg("Artemis: production auth procedure starting")
		temporalAuthCfg = temporalProdAuthConfig
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		inMemSecrets, sw := auth_startup.RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		temporalAuthCfg = temporalProdAuthConfig
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		inMemSecrets, sw := auth_startup.RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		temporalAuthCfg = tc.ProdLocalTemporalAuthArtemis
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
	case "local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
		temporalAuthCfg = tc.ProdLocalTemporalAuthArtemis
		artemis_network_cfgs.InitArtemisLocalTestConfigs()
	}

	log.Info().Msgf("Artemis %s temporal auth and init procedure starting", env)
	artemis_ethereum_transcations.InitEthereumBroadcasters(ctx, temporalAuthCfg)
	log.Info().Msgf("Artemis %s temporal auth and init procedure succeeded", env)
}
