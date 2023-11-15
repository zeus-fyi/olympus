package tyche_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/configs"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_mev_tx_fetcher "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/mev"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
	artemis_trade_executor "github.com/zeus-fyi/olympus/pkg/artemis/trading/executor"
	"github.com/zeus-fyi/olympus/pkg/artemis/trading/pricing/price_quoter"
	"github.com/zeus-fyi/olympus/pkg/athena"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	tyche_metrics "github.com/zeus-fyi/olympus/tyche/metrics"
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
		log.Info().Msg("Tyche: production auth procedure starting")
		temporalAuthCfg = temporalProdAuthConfig
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		inMemSecrets, sw := auth_startup.RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
		dynamoDBCreds.AccessKey = sw.AccessKeyHydraDynamoDB
		dynamoDBCreds.AccessSecret = sw.SecretKeyHydraDynamoDB
		artemis_orchestration_auth.Bearer = sw.BearerToken
		price_quoter.ZeroXApiKey = sw.ZeroXApiKey
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
		//artemis_trading_cache.InitProductionRedis(ctx)
		artemis_trading_cache.InitBackupProductionRedis(ctx)
		iris_redis.InitProductionBackupRedisIrisCache(ctx)
		//iris_redis.InitProductionRedisIrisCache(ctx)
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.ProdLocalAuthKeysCfg)
		temporalAuthCfg = tc.DevTemporalAuth
		inMemSecrets, sw := auth_startup.RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		dynamoDBCreds.AccessKey = sw.AccessKeyHydraDynamoDB
		dynamoDBCreds.AccessSecret = sw.SecretKeyHydraDynamoDB
		price_quoter.ZeroXApiKey = sw.ZeroXApiKey
		artemis_orchestration_auth.Bearer = sw.BearerToken
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
		iris_redis.InitLocalTestProductionRedisIrisCache(ctx)
	case "local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
		temporalAuthCfg = tc.DevTemporalAuth
		dynamoDBCreds.AccessKey = tc.AwsAccessKeyDynamoDB
		dynamoDBCreds.AccessSecret = tc.AwsSecretKeyDynamoDB
		price_quoter.ZeroXApiKey = tc.ZeroXApiKey
		artemis_orchestration_auth.Bearer = tc.ProductionLocalTemporalBearerToken
		authKeysCfg = tc.DevAuthKeysCfg
		artemis_network_cfgs.InitArtemisLocalTestConfigs()
		iris_redis.InitLocalTestRedisIrisCache(ctx)
	}
	dynamoDBCreds.Region = "us-west-1"
	log.Info().Msg("Tyche: DynamoDB connection starting")
	artemis_orchestration_auth.InitMevDynamoDBClient(dynamoDBCreds)
	log.Info().Msg("Tyche: DynamoDB connection succeeded")

	log.Info().Msg("Tyche: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Tyche: PG connection succeeded")

	log.Info().Msg("Tyche: DigitalOceanS3AuthClient starting")
	athena.AthenaS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg)

	log.Info().Msg("Tyche: InitTokenFilter starting")
	artemis_trading_cache.InitTokenFilter(ctx)
	log.Info().Msg("Tyche: InitTokenFilter succeeded")

	log.Info().Msg("Tyche: InitFlashbots starting")
	age := encryption.NewAge(authKeysCfg.AgePrivKey, authKeysCfg.AgePubKey)
	tyche_metrics.InitTycheMetrics(ctx)
	artemis_trading_cache.InitWeb3Client()
	artemis_trade_executor.InitMainnetAuxiliaryTradingUtils(ctx, age)
	artemis_trade_executor.InitGoerliAuxiliaryTradingUtils(ctx, age)
	log.Info().Msg("Tyche: InitFlashbots succeeded")

	artemis_mev_tx_fetcher.InitTycheUniswap(ctx, artemis_orchestration_auth.Bearer)
	log.Info().Msg("Tyche: InitArtemisUniswap succeeded")

}
