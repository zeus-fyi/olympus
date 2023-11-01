package artemis_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/configs"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_mev_tx_fetcher "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/mev"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	artemis_ethereum_transcations "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/transcations"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trade_executor "github.com/zeus-fyi/olympus/pkg/artemis/trading/executor"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
	"github.com/zeus-fyi/olympus/pkg/athena"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
)

var (
	temporalProdAuthConfig = temporal_auth.TemporalAuth{
		ClientCertPath:   "/etc/ssl/certs/ca.pem",
		ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
		Namespace:        "production-artemis.ngb72",
		HostPort:         "production-artemis.ngb72.tmprl.cloud:7233",
	}
	dynamoDBCreds = dynamodb_client.DynamoDBCredentials{}
)

func SetConfigByEnv(ctx context.Context, env string) {
	switch env {
	case "production":
		log.Info().Msg("Artemis: production auth procedure starting")
		temporalAuthCfg = temporalProdAuthConfig
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		inMemSecrets, sw := auth_startup.RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
		dynamoDBCreds.AccessKey = sw.AccessKeyHydraDynamoDB
		dynamoDBCreds.AccessSecret = sw.SecretKeyHydraDynamoDB
		artemis_trading_cache.InitProductionRedis(ctx)
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.ProdLocalAuthKeysCfg)
		temporalAuthCfg = tc.DevTemporalAuth
		inMemSecrets, sw := auth_startup.RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		dynamoDBCreds.AccessKey = sw.AccessKeyHydraDynamoDB
		dynamoDBCreds.AccessSecret = sw.SecretKeyHydraDynamoDB
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
	case "local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
		temporalAuthCfg = tc.DevTemporalAuth
		dynamoDBCreds.AccessKey = tc.AwsAccessKeyDynamoDB
		dynamoDBCreds.AccessSecret = tc.AwsSecretKeyDynamoDB
		authKeysCfg = tc.DevAuthKeysCfg
		artemis_network_cfgs.InitArtemisLocalTestConfigs()
	}
	dynamoDBCreds.Region = "us-west-1"

	log.Info().Msg("Artemis: DynamoDB connection starting")
	artemis_orchestration_auth.InitMevDynamoDBClient(dynamoDBCreds)
	log.Info().Msg("Artemis: DynamoDB connection succeeded")

	log.Info().Msg("Artemis: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Artemis: PG connection succeeded")

	log.Info().Msg("Artemis: InitTokenFilter starting")
	artemis_trading_cache.InitTokenFilter(ctx)
	log.Info().Msg("Artemis: InitTokenFilter succeeded")

	log.Info().Msgf("Artemis %s orchestration retrieving auth token", env)
	artemis_orchestration_auth.Bearer = auth_startup.FetchTemporalAuthBearer(ctx)
	artemis_test_cache.LiveTestNetwork.AddBearerToken(artemis_orchestration_auth.Bearer)
	log.Info().Msgf("Artemis %s orchestration retrieving auth token done", env)

	log.Info().Msgf("Artemis InitEthereumBroadcasters: %s temporal auth and init procedure starting", env)
	artemis_ethereum_transcations.InitEthereumBroadcasters(ctx, temporalAuthCfg)
	log.Info().Msgf("Artemis InitEthereumBroadcasters: %s temporal auth and init procedure succeeded", env)

	log.Info().Msgf("Artemis InitMevWorkers: %s temporal auth and init procedure starting", env)
	artemis_mev_tx_fetcher.InitMevWorkers(ctx, temporalAuthCfg)
	log.Info().Msgf("Artemis InitMevWorkers: %s temporal auth and init procedure succeeded", env)

	log.Info().Msgf("Artemis %s init flashbots client", env)
	artemis_trading_cache.InitWeb3Client()
	athena.AthenaS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg)
	age := encryption.NewAge(authKeysCfg.AgePrivKey, authKeysCfg.AgePubKey)
	artemis_trade_executor.InitMainnetAuxiliaryTradingUtils(ctx, age)
	artemis_trade_executor.InitGoerliAuxiliaryTradingUtils(ctx, age)
	log.Info().Msgf("Artemis %s done init flashbots client", env)
}
