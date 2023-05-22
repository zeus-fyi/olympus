package artemis_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_mev_tx_fetcher "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/mev"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
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
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.ProdLocalAuthKeysCfg)
		temporalAuthCfg = tc.DevTemporalAuth
		inMemSecrets, sw := auth_startup.RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
	case "local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
		temporalAuthCfg = tc.DevTemporalAuth
		artemis_network_cfgs.InitArtemisLocalTestConfigs()
	}

	log.Info().Msg("Artemis: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Artemis: PG connection succeeded")

	log.Info().Msgf("Artemis %s orchestration retrieving auth token", env)
	artemis_orchestration_auth.Bearer = auth_startup.FetchTemporalAuthBearer(ctx)
	log.Info().Msgf("Artemis %s orchestration retrieving auth token done", env)

	log.Info().Msgf("Artemis InitEthereumBroadcasters: %s temporal auth and init procedure starting", env)
	artemis_ethereum_transcations.InitEthereumBroadcasters(ctx, temporalAuthCfg)
	log.Info().Msgf("Artemis InitEthereumBroadcasters: %s temporal auth and init procedure succeeded", env)

	log.Info().Msgf("Artemis InitMevWorkers: %s temporal auth and init procedure starting", env)
	artemis_mev_tx_fetcher.InitMevWorkers(ctx, temporalAuthCfg)
	log.Info().Msgf("Artemis InitMevWorkers: %s temporal auth and init procedure succeeded", env)
}
