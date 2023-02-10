package hydra_server

import (
	"context"
	"errors"
	dynamodb_client "github.com/zeus-fyi/olympus/datastores/dynamodb"
	ethereum_slashing_protection_watermarking "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/slashing_protection"
	"github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/dynamodb_web3signer"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
)

var (
	temporalProdAuthConfig = temporal_auth.TemporalAuth{
		ClientCertPath:   "/etc/ssl/certs/ca.pem",
		ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
		Namespace:        "production-artemis.ngb72",
		HostPort:         "production-artemis.ngb72.tmprl.cloud:7233",
	}
	dynamoCreds = dynamodb_client.DynamoDBCredentials{
		Region: "us-west-1",
	}
)

func SetConfigByEnv(ctx context.Context, env string) {
	switch env {
	case "production":
		log.Info().Msg("Artemis: production auth procedure starting")
		temporalAuthCfg = temporalProdAuthConfig
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		inMemSecrets, sw := auth_startup.RunHydraDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
		dynamoCreds.AccessKey = sw.AccessKeyHydraDynamoDB
		dynamoCreds.AccessSecret = sw.SecretKeyHydraDynamoDB
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		temporalAuthCfg = temporalProdAuthConfig
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		inMemSecrets, sw := auth_startup.RunHydraDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		dynamoCreds.AccessKey = sw.AccessKeyHydraDynamoDB
		dynamoCreds.AccessSecret = sw.SecretKeyHydraDynamoDB
		temporalAuthCfg = tc.ProdLocalTemporalAuthArtemis
		auth_startup.InitArtemisEthereum(ctx, inMemSecrets, sw)
	case "local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
		temporalAuthCfg = tc.ProdLocalTemporalAuthArtemis
		dynamoCreds.AccessKey = tc.AwsAccessKey
		dynamoCreds.AccessSecret = tc.AwsSecretKey
		artemis_network_cfgs.InitArtemisLocalTestConfigs()
	}
	log.Info().Msg("Hydra: InitPG connecting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Hydra: InitPG done")

	switch Workload.ProtocolNetworkID {
	case hestia_req_types.EthereumEphemeryProtocolNetworkID:
		log.Info().Msg("Hydra: ProtocolNetworkID (ephemery)")
		ethereum_slashing_protection_watermarking.Network = "ephemery"
	case hestia_req_types.EthereumMainnetProtocolNetworkID:
		log.Info().Msg("Hydra: ProtocolNetworkID (mainnet)")
		ethereum_slashing_protection_watermarking.Network = "mainnet"
	default:
		err := errors.New("invalid or unsupported protocol network id")
		log.Ctx(ctx).Err(err)
		panic(err)
	}
	log.Info().Msg("Hydra: InitDynamoDB connecting")
	dynamodb_web3signer_client.InitWeb3SignerDynamoDBClient(ctx, dynamoCreds)
	log.Info().Msg("Hydra: InitDynamoDB done")

	log.Info().Msgf("Hydra %s artemis orchestration retrieving auth token", env)
	artemis_orchestration_auth.Bearer = auth_startup.FetchTemporalAuthBearer(ctx)
	log.Info().Msgf("Hydra %s artemis orchestration retrieving auth token done", env)
}
