package tyche_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/configs"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
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
		price_quoter.ZeroXApiKey = sw.ZeroXApiKey
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.ProdLocalAuthKeysCfg)
		temporalAuthCfg = tc.DevTemporalAuth
		inMemSecrets, sw := auth_startup.RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		dynamoDBCreds.AccessKey = sw.AccessKeyHydraDynamoDB
		dynamoDBCreds.AccessSecret = sw.SecretKeyHydraDynamoDB
		price_quoter.ZeroXApiKey = sw.ZeroXApiKey
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
	case "local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
		temporalAuthCfg = tc.DevTemporalAuth
		dynamoDBCreds.AccessKey = tc.AwsAccessKeyDynamoDB
		dynamoDBCreds.AccessSecret = tc.AwsSecretKeyDynamoDB
		price_quoter.ZeroXApiKey = tc.ZeroXApiKey
		authKeysCfg = tc.DevAuthKeysCfg
		artemis_network_cfgs.InitArtemisLocalTestConfigs()
	}
	dynamoDBCreds.Region = "us-west-1"
	log.Info().Msg("Tyche: DynamoDB connection starting")
	artemis_orchestration_auth.InitMevDynamoDBClient(dynamoDBCreds)
	log.Info().Msg("Tyche: DynamoDB connection succeeded")

	log.Info().Msg("Tyche: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Tyche: PG connection succeeded")

	log.Info().Msg("Athena: DigitalOceanS3AuthClient starting")
	athena.AthenaS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg)

	log.Info().Msg("Tyche: InitTokenFilter starting")
	artemis_trading_cache.InitTokenFilter(ctx)
	log.Info().Msg("Tyche: InitTokenFilter succeeded")

	log.Info().Msg("Tyche: InitFlashbots starting")
	age := encryption.NewAge(authKeysCfg.AgePrivKey, authKeysCfg.AgePubKey)
	tm := tyche_metrics.InitTycheMetrics(ctx)
	artemis_trade_executor.InitMainnetAuxiliaryTradingUtils(ctx, age, &tm)
	artemis_trade_executor.InitGoerliAuxiliaryTradingUtils(ctx, age)
	log.Info().Msg("Tyche: InitFlashbots succeeded")

}
