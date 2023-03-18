package hestia_server

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	v1hestia "github.com/zeus-fyi/olympus/hestia/api/v1"
	hestia_web_router "github.com/zeus-fyi/olympus/hestia/web"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	eth_validators_service_requests "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validators_service_requests"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	artemis_client "github.com/zeus-fyi/zeus/pkg/artemis/client"
)

var (
	cfg                = Config{}
	authKeysCfg        auth_keys_config.AuthKeysCfg
	env                string
	temporalAuthConfig = temporal_auth.TemporalAuth{
		ClientCertPath:   "/etc/ssl/certs/ca.pem",
		ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
		Namespace:        "production-artemis.ngb72",
		HostPort:         "production-artemis.ngb72.tmprl.cloud:7233",
	}
	awsRegion  = "us-west-1"
	awsAuthCfg = aegis_aws_auth.AuthAWS{
		Region:    awsRegion,
		AccessKey: "",
		SecretKey: "",
	}
)

func Hestia() {
	cfg.Host = "0.0.0.0"
	srv := NewHestiaServer(cfg)
	// Echo instance
	ctx := context.Background()
	switch env {
	case "production":
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunHestiaDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
		awsAuthCfg = sw.SecretsManagerAuthAWS
		awsAuthCfg.Region = awsRegion
		artemis_validator_service_groups_models.ArtemisClient = artemis_client.NewDefaultArtemisClient(sw.BearerToken)
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		temporalAuthConfig = tc.ProdLocalTemporalAuthArtemis
		awsAuthCfg.AccessKey = tc.AwsAccessKeySecretManager
		awsAuthCfg.SecretKey = tc.AwsSecretKeySecretManager
		artemis_validator_service_groups_models.ArtemisClient = artemis_client.NewDefaultArtemisClient(tc.ProductionLocalTemporalBearerToken)
	case "local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
		temporalAuthConfig = tc.ProdLocalTemporalAuthArtemis
		awsAuthCfg.AccessKey = tc.AwsAccessKeySecretManager
		awsAuthCfg.SecretKey = tc.AwsSecretKeySecretManager
		artemis_validator_service_groups_models.ArtemisClient = artemis_client.NewDefaultArtemisClient(tc.ProductionLocalTemporalBearerToken)
	}
	log.Info().Msg("Hestia: AWS Secrets Manager connection starting")
	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, awsAuthCfg)
	log.Info().Msg("Hestia: AWS Secrets Manager connected")

	log.Info().Msg("Hestia: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Hestia: PG connection connected")

	log.Info().Msgf("Hestia %s artemis orchestration retrieving auth token", env)
	artemis_orchestration_auth.Bearer = auth_startup.FetchTemporalAuthBearer(ctx)
	log.Info().Msgf("Hestia %s artemis orchestration retrieving auth token done", env)

	log.Info().Msgf("Hestia %s artemis orchestration setting up zeus client", env)
	eth_validators_service_requests.InitZeusClientValidatorServiceGroup(ctx)
	log.Info().Msgf("Hestia %s artemis orchestration zeus client setup complete", env)

	log.Info().Msg("Hestia: InitArtemisEthereumEphemeryValidatorsRequestsWorker Starting")
	// NOTE: inits at least one worker, then reuses the connection
	// ephemery
	eth_validators_service_requests.InitArtemisEthereumEphemeryValidatorsRequestsWorker(ctx, temporalAuthConfig)
	// connect
	c := eth_validators_service_requests.ArtemisEthereumEphemeryValidatorsRequestsWorker.ConnectTemporalClient()
	defer c.Close()

	log.Info().Msg("Hestia: InitArtemisEthereumEphemeryValidatorsRequestsWorker Done")
	eth_validators_service_requests.ArtemisEthereumEphemeryValidatorsRequestsWorker.Worker.RegisterWorker(c)
	err := eth_validators_service_requests.ArtemisEthereumEphemeryValidatorsRequestsWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Hestia: %s ArtemisEthereumEphemeryValidatorsRequestsWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Hestia: InitArtemisEthereumGoerliValidatorsRequestsWorker")
	eth_validators_service_requests.InitArtemisEthereumGoerliValidatorsRequestsWorker(ctx, temporalAuthConfig)
	log.Info().Msg("Hestia: Starting InitArtemisEthereumGoerliValidatorsRequestsWorker")
	eth_validators_service_requests.ArtemisEthereumGoerliValidatorsRequestsWorker.Worker.RegisterWorker(c)
	err = eth_validators_service_requests.ArtemisEthereumGoerliValidatorsRequestsWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Hestia: %s ArtemisEthereumGoerliValidatorsRequestsWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}

	log.Info().Msg("Hestia: InitArtemisEthereumMainnetValidatorsRequestsWorker")
	// mainnet
	eth_validators_service_requests.InitArtemisEthereumMainnetValidatorsRequestsWorker(ctx, temporalAuthConfig)
	log.Info().Msg("Hestia: Starting InitArtemisEthereumMainnetValidatorsRequestsWorker")
	eth_validators_service_requests.ArtemisEthereumMainnetValidatorsRequestsWorker.Worker.RegisterWorker(c)
	err = eth_validators_service_requests.ArtemisEthereumMainnetValidatorsRequestsWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Hestia: %s ArtemisEthereumMainnetValidatorsRequestsWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Hestia: InitArtemisEthereumMainnetValidatorsRequestsWorker Done")
	log.Info().Msg("Hestia: InitArtemisEthereumMainnetValidatorsRequestsWorker Starting Server")

	if env == "local" || env == "production-local" {
		srv.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			AllowCredentials: true,
		}))
	} else {
		srv.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"https://cloud.zeus.fyi"},
			AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			AllowCredentials: true,
		}))
	}
	srv.E = v1hestia.Routes(srv.E)
	srv.E = hestia_web_router.WebRoutes(srv.E)
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9002", "server port")
	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")
	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Storing Internal Data",
	Short: "A microservice for internal configurations",
	Run: func(cmd *cobra.Command, args []string) {
		Hestia()
	},
}
