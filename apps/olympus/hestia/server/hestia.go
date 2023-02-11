package hestia_server

import (
	"context"
	hestia_aws_secrets_auth "github.com/zeus-fyi/olympus/hestia/auth"
	aws_secrets "github.com/zeus-fyi/zeus/pkg/aegis/aws"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	v1hestia "github.com/zeus-fyi/olympus/hestia/api/v1"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	eth_validators_service_requests "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validators_service_requests"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

var (
	cfg                    = Config{}
	authKeysCfg            auth_keys_config.AuthKeysCfg
	env                    string
	temporalProdAuthConfig = temporal_auth.TemporalAuth{
		ClientCertPath:   "/etc/ssl/certs/ca.pem",
		ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
		Namespace:        "production-artemis.ngb72",
		HostPort:         "production-artemis.ngb72.tmprl.cloud:7233",
	}
	awsAuthCfg = aws_secrets.AuthAWS{
		Region:    "us-west-1",
		AccessKey: "",
		SecretKey: "",
	}
)

func Hestia() {
	cfg.Host = "0.0.0.0"
	srv := NewHestiaServer(cfg)
	// Echo instance
	srv.E = v1hestia.Routes(srv.E)
	ctx := context.Background()
	switch env {
	case "production":
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunHestiaDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
		awsAuthCfg = sw.SecretsManagerAuthAWS
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		temporalProdAuthConfig = tc.ProdLocalTemporalAuthArtemis
		awsAuthCfg.AccessKey = tc.AwsAccessKeySecretManager
		awsAuthCfg.SecretKey = tc.AwsSecretKeySecretManager
	case "local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
		temporalProdAuthConfig = tc.ProdLocalTemporalAuthArtemis
		awsAuthCfg.AccessKey = tc.AwsAccessKeySecretManager
		awsAuthCfg.SecretKey = tc.AwsSecretKeySecretManager
	}
	log.Info().Msg("Hestia: AWS Secrets Manager connection starting")
	hestia_aws_secrets_auth.InitHestiaSecretManagerAuthAWS(ctx, awsAuthCfg)
	log.Info().Msg("Hestia: AWS Secrets Manager connected")

	log.Info().Msg("Hestia: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Hestia: PG connection connected")
	// NOTE: inits at least one worker, then reuses the connection
	// ephemery
	eth_validators_service_requests.InitArtemisEthereumEphemeryValidatorsRequestsWorker(ctx, temporalProdAuthConfig)
	// connect
	c := eth_validators_service_requests.ArtemisEthereumEphemeryValidatorsRequestsWorker.ConnectTemporalClient()
	defer c.Close()

	log.Info().Msg("Hestia: Starting InitArtemisEthereumEphemeryValidatorsRequestsWorker")
	eth_validators_service_requests.ArtemisEthereumEphemeryValidatorsRequestsWorker.Worker.RegisterWorker(c)
	err := eth_validators_service_requests.ArtemisEthereumEphemeryValidatorsRequestsWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Hestia: %s ArtemisEthereumEphemeryValidatorsRequestsWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}

	// mainnet
	eth_validators_service_requests.InitArtemisEthereumMainnetValidatorsRequestsWorker(ctx, temporalProdAuthConfig)
	log.Info().Msg("Hestia: Starting InitArtemisEthereumMainnetValidatorsRequestsWorker")
	eth_validators_service_requests.ArtemisEthereumMainnetValidatorsRequestsWorker.Worker.RegisterWorker(c)
	err = eth_validators_service_requests.ArtemisEthereumMainnetValidatorsRequestsWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Hestia: %s ArtemisEthereumMainnetValidatorsRequestsWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}

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
