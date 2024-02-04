package iris_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/configs"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	v1_iris "github.com/zeus-fyi/olympus/iris/api/v1"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	proxy_anvil "github.com/zeus-fyi/olympus/pkg/iris/proxy/anvil"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	iris_serverless "github.com/zeus-fyi/olympus/pkg/iris/serverless"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

var (
	temporalProdAuthConfig = temporal_auth.TemporalAuth{
		ClientCertPath:   "/etc/ssl/certs/ca.pem",
		ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
		Namespace:        "production-iris.ngb72",
		HostPort:         "production-iris.ngb72.tmprl.cloud:7233",
	}
	dynamoDBCreds = dynamodb_client.DynamoDBCredentials{}
	awsRegion     = "us-west-1"
	awsAuthCfg    = aegis_aws_auth.AuthAWS{
		Region:    awsRegion,
		AccessKey: "",
		SecretKey: "",
	}
)

func SetConfigByEnv(ctx context.Context, env string) {
	v1_iris.Env = env
	switch env {
	case "production":
		log.Info().Msg("Iris: production auth procedure starting")
		temporalAuthCfg = temporalProdAuthConfig
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		inMemSecrets, sw := auth_startup.RunIrisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
		dynamoDBCreds.AccessKey = sw.AccessKeyHydraDynamoDB
		dynamoDBCreds.AccessSecret = sw.SecretKeyHydraDynamoDB
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
		log.Info().Msg("Iris: AWS Secrets Manager connection starting")
		awsAuthCfg = sw.SecretsManagerAuthAWS
		awsAuthCfg.Region = awsRegion
		log.Info().Msg("Iris: AWS Secrets Manager connected")
		iris_redis.InitProductionBackupRedisIrisCache(ctx)
		//iris_redis.InitProductionRedisIrisCache(ctx)
		artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, awsAuthCfg)
	case "production-local":
		auth_startup.Ksp.DirIn = "../configs"
		auth_startup.Sp.DirIn = "../configs"
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.ProdLocalAuthKeysCfg)
		temporalAuthCfg = tc.DevTemporalAuth
		inMemSecrets, sw := auth_startup.RunIrisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		dynamoDBCreds.AccessKey = sw.AccessKeyHydraDynamoDB
		dynamoDBCreds.AccessSecret = sw.SecretKeyHydraDynamoDB
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
		iris_redis.InitLocalTestProductionRedisIrisCache(ctx)

		awsAuthCfg.AccessKey = tc.AwsAccessKeySecretManager
		awsAuthCfg.SecretKey = tc.AwsSecretKeySecretManager
		artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, awsAuthCfg)
	case "local":
		auth_startup.Ksp.DirIn = "../configs"
		auth_startup.Sp.DirIn = "../configs"
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
		temporalAuthCfg = tc.DevTemporalAuth
		dynamoDBCreds.AccessKey = tc.AwsAccessKeyDynamoDB
		dynamoDBCreds.AccessSecret = tc.AwsSecretKeyDynamoDB
		artemis_network_cfgs.InitArtemisLocalTestConfigs()
		iris_redis.InitLocalTestRedisIrisCache(ctx)
		awsAuthCfg.AccessKey = tc.AwsAccessKeySecretManager
		awsAuthCfg.SecretKey = tc.AwsSecretKeySecretManager
		artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, awsAuthCfg)

	default:
		iris_redis.InitLocalTestRedisIrisCache(ctx)
	}

	//dynamoDBCreds.Region = "us-west-1"
	//
	//log.Info().Msg("Artemis: DynamoDB connection starting")
	//artemis_orchestration_auth.InitMevDynamoDBClient(dynamoDBCreds)
	//log.Info().Msg("Artemis: DynamoDB connection succeeded")

	log.Info().Msg("Iris: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Iris: PG connection succeeded")

	log.Info().Msgf("Iris %s orchestration retrieving auth token", env)
	artemis_orchestration_auth.Bearer = auth_startup.FetchTemporalAuthBearer(ctx)
	log.Info().Msgf("Iris %s orchestration retrieving auth token done", env)

	log.Info().Msgf("Iris InitIrisApiRequestsWorker: %s temporal auth and init procedure starting", env)
	iris_api_requests.InitIrisApiRequestsWorker(ctx, temporalAuthCfg)
	log.Info().Msgf("Iris InitIrisApiRequestsWorker: %s temporal auth and init procedure succeeded", env)

	log.Info().Msgf("Iris InitIrisCacheWorker: %s temporal auth and init procedure starting", env)
	iris_api_requests.InitIrisCacheWorker(ctx, temporalAuthCfg)
	log.Info().Msgf("Iris InitIrisCacheWorker: %s temporal auth and init procedure succeeded", env)

	log.Info().Msgf("Iris InitIrisPlatformServicesWorker: %s temporal auth and init procedure starting", env)
	iris_serverless.InitIrisPlatformServicesWorker(ctx, temporalAuthCfg)
	log.Info().Msgf("Iris InitIrisPlatformServicesWorker: %s temporal auth and init procedure succeeded", env)

	proxy_anvil.InitAnvilProxy()
}
