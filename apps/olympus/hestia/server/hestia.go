package hestia_server

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	v1hestia "github.com/zeus-fyi/olympus/hestia/api/v1"
	v1_ethereum_aws "github.com/zeus-fyi/olympus/hestia/api/v1/ethereum/aws"
	hestia_quiknode_v1_routes "github.com/zeus-fyi/olympus/hestia/api/v1/quiknode"
	hestia_web_router "github.com/zeus-fyi/olympus/hestia/web"
	hestia_billing "github.com/zeus-fyi/olympus/hestia/web/billing"
	hestia_iris_dashboard "github.com/zeus-fyi/olympus/hestia/web/iris"
	hestia_login "github.com/zeus-fyi/olympus/hestia/web/login"
	hestia_mev "github.com/zeus-fyi/olympus/hestia/web/mev"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	artemis_validator_signature_service_routing "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/signature_routing"
	eth_validators_service_requests "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validators_service_requests"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	hermes_email_notifications "github.com/zeus-fyi/olympus/pkg/hermes/email"
	platform_service_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/iris/orchestrations"
	quicknode_orchestrations "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode/orchestrations"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	artemis_client "github.com/zeus-fyi/zeus/pkg/artemis/client"
	"golang.org/x/oauth2"
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
	temporalAuthConfigHestia = temporal_auth.TemporalAuth{
		ClientCertPath:   "/etc/ssl/certs/ca.pem",
		ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
		Namespace:        "production-hestia.ngb72",
		HostPort:         "production-hestia.ngb72.tmprl.cloud:7233",
	}
	temporalAuthConfigKronos = temporal_auth.TemporalAuth{
		ClientCertPath:   "/etc/ssl/certs/ca.pem",
		ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
		Namespace:        "kronos.ngb72",
		HostPort:         "kronos.ngb72.tmprl.cloud:7233",
	}
	awsRegion  = "us-west-1"
	awsAuthCfg = aegis_aws_auth.AuthAWS{
		Region:    awsRegion,
		AccessKey: "",
		SecretKey: "",
	}
	awsSESAuthCfg = aegis_aws_auth.AuthAWS{
		Region:    awsRegion,
		AccessKey: "",
		SecretKey: "",
	}
)

const (
	SelectedRouteHeader        = "X-Selected-Route"
	SelectedLatencyHeader      = "X-Response-Latency-Milliseconds"
	SelectedRouteGroupHeader   = "X-Route-Group"
	SelectedResponseReceivedAt = "X-Response-Received-At-UTC"
	AdaptiveMetricsKey         = "X-Adaptive-Metrics-Key"
)

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

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
		sw.SESAuthAWS.Region = awsRegion
		if len(sw.TwitterMbClientID) <= 0 || len(sw.TwitterMbClientSecret) <= 0 {
			log.Warn().Msg("Hestia: TwitterClientID or TwitterClientSecret is empty")
		}

		authorizeURL := "https://twitter.com/i/oauth2/authorize"
		tokenURL := "https://api.twitter.com/2/oauth2/token"
		conf := &oauth2.Config{

			// "https://hestia.zeus.fyi/social/v1/auth/twitter/callback"
			RedirectURL:  "https://cloud.zeus.fyi/social/v1/twitter/callback",
			ClientID:     sw.TwitterMbClientID,
			ClientSecret: sw.TwitterMbClientSecret,
			Scopes: []string{
				"bookmark.write",
				"bookmark.read",
				"tweet.read",
				"tweet.write",
				"tweet.moderate.write",
				"users.read",
				"follows.read",
				"follows.write",
				"offline.access",
				"space.read",
				"mute.read",
				"mute.write",
				"like.read",
				"like.write",
				"list.read",
				"list.write",
				"block.read",
				"block.write",
			}, Endpoint: oauth2.Endpoint{
				AuthURL:  authorizeURL,
				TokenURL: tokenURL,
			},
		}
		hestia_login.TwitterOAuthConfig = conf
		hestia_iris_dashboard.JWTAuthSecret = sw.QuickNodeJWT
		hestia_quiknode_v1_routes.QuickNodePassword = sw.QuickNodePassword
		if len(hestia_quiknode_v1_routes.QuickNodePassword) <= 0 {
			log.Fatal().Msg("Hestia: QuickNodePassword is empty")
			misc.DelayedPanic(errors.New("hestia: QuickNodePassword is empty"))
		}
		if len(hestia_iris_dashboard.JWTAuthSecret) <= 0 {
			log.Fatal().Msg("Hestia: JWTAuthSecret is empty")
			misc.DelayedPanic(errors.New("hestia: QuickNodePassword is empty"))
		}
		artemis_validator_service_groups_models.ArtemisClient = artemis_client.NewDefaultArtemisClient(sw.BearerToken)
		artemis_orchestration_auth.Bearer = sw.BearerToken
		//hermes_email_notifications.Hermes = hermes_email_notifications.InitHermesSESEmailNotifications(ctx, sw.SESAuthAWS)
		hermes_email_notifications.InitHermesSendGridClient(ctx, sw.SendGridAPIKey)
		hestia_stripe.InitStripe(sw.StripeSecretKey)
		kronos_helix.InitPagerDutyAlertClient(sw.PagerDutyApiKey)
		kronos_helix.PdAlertGenericWfIssuesEvent.RoutingKey = sw.PagerDutyRoutingKey
		if len(kronos_helix.PdAlertGenericWfIssuesEvent.RoutingKey) <= 0 {
			log.Fatal().Msg("Hestia: PagerDutyRoutingKey is empty")
			misc.DelayedPanic(errors.New("hestia: PagerDutyRoutingKey is empty"))
		}
		hestia_login.GoogleOAuthConfig.ClientID = sw.GoogClientID
		hestia_login.GoogleOAuthConfig.ClientSecret = sw.GoogClientSecret
		//hestia_analytics.GtagApiSecret = sw.GoogGtagSecret
		hestia_login.DiscordRedirectURI = "https://hestia.zeus.fyi/discord/callback"

		hestia_login.DiscordClientID = sw.DiscordAuthConfig.DiscordClientID
		hestia_login.DiscordClientSecret = sw.DiscordAuthConfig.DiscordClientSecret
		hestia_login.SetConf(sw.DiscordAuthConfig.DiscordClientID, sw.DiscordAuthConfig.DiscordClientSecret)
		//DiscordRedirectURI
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		temporalAuthConfig = tc.DevTemporalAuth
		temporalAuthConfigHestia = tc.DevTemporalAuth
		temporalAuthConfigKronos = tc.DevTemporalAuth
		awsAuthCfg.AccessKey = tc.AwsAccessKeySecretManager
		awsAuthCfg.SecretKey = tc.AwsSecretKeySecretManager

		hestia_iris_dashboard.JWTAuthSecret = tc.QuickNodeMarketplace.JWTToken
		hestia_quiknode_v1_routes.QuickNodePassword = tc.QuickNodeMarketplace.Password

		artemis_validator_service_groups_models.ArtemisClient = artemis_client.NewDefaultArtemisClient(tc.ProductionLocalTemporalBearerToken)
		awsSESAuthCfg.AccessKey = tc.AwsAccessKeySES
		awsSESAuthCfg.SecretKey = tc.AwsSecretKeySES
		//hermes_email_notifications.Hermes = hermes_email_notifications.InitHermesSESEmailNotifications(ctx, awsSESAuthCfg)
		hermes_email_notifications.InitHermesSendGridClient(ctx, tc.SendGridAPIKey)
		hestia_stripe.InitStripe(tc.StripeTestSecretAPIKey)
		kronos_helix.InitPagerDutyAlertClient(tc.PagerDutyApiKey)
		kronos_helix.PdAlertGenericWfIssuesEvent.RoutingKey = tc.PagerDutyRoutingKey
		platform_service_orchestrations.IrisApiUrl = "http://localhost:8080"
		quicknode_orchestrations.IrisApiUrl = "http://localhost:8080"
		hestia_login.GoogleOAuthConfig.ClientID = tc.GoogClientID
		hestia_login.GoogleOAuthConfig.ClientSecret = tc.GoogClientSecret
		//hestia_analytics.GtagApiSecret = tc.GoogTagSecret
		hera_openai.InitHeraOpenAI(tc.OpenAIAuth)
		hestia_mev.PromqlProxy = "http://localhost:8000"
		hestia_login.DiscordRedirectURI = "http://localhost:9002/discord/callback"
		hestia_login.DiscordClientID = tc.DiscordClientID
		hestia_login.DiscordClientSecret = tc.DiscordClientSecret
		hestia_login.SetConf(tc.DiscordClientID, tc.DiscordClientSecret)
	case "local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
		temporalAuthConfig = tc.DevTemporalAuth
		temporalAuthConfigHestia = tc.DevTemporalAuth
		temporalAuthConfigKronos = tc.DevTemporalAuth
		awsAuthCfg.AccessKey = tc.AwsAccessKeySecretManager
		awsAuthCfg.SecretKey = tc.AwsSecretKeySecretManager
		hestia_iris_dashboard.JWTAuthSecret = tc.QuickNodeMarketplace.JWTToken
		hestia_quiknode_v1_routes.QuickNodePassword = tc.QuickNodeMarketplace.Password
		artemis_validator_service_groups_models.ArtemisClient = artemis_client.NewDefaultArtemisClient(tc.ProductionLocalTemporalBearerToken)
		awsSESAuthCfg.AccessKey = tc.AwsAccessKeySES
		awsSESAuthCfg.SecretKey = tc.AwsSecretKeySES
		hermes_email_notifications.Hermes = hermes_email_notifications.InitHermesSESEmailNotifications(ctx, awsSESAuthCfg)
		hermes_email_notifications.InitHermesSendGridClient(ctx, tc.SendGridAPIKey)
		hestia_stripe.InitStripe(tc.StripeTestSecretAPIKey)
		kronos_helix.InitPagerDutyAlertClient(tc.PagerDutyApiKey)
		kronos_helix.PdAlertGenericWfIssuesEvent.RoutingKey = tc.PagerDutyRoutingKey
		platform_service_orchestrations.IrisApiUrl = "http://localhost:8080"
		quicknode_orchestrations.IrisApiUrl = "http://localhost:8080"
		hestia_login.GoogleOAuthConfig.ClientID = tc.GoogClientID
		hestia_login.GoogleOAuthConfig.ClientSecret = tc.GoogClientSecret
		//hestia_analytics.GtagApiSecret = tc.GoogTagSecret
		hestia_mev.PromqlProxy = "http://localhost:8000"
		hera_openai.InitHeraOpenAI(tc.OpenAIAuth)
		hestia_login.DiscordRedirectURI = "http://localhost:9002/discord/callback"
		hestia_login.DiscordClientID = tc.DiscordClientID
		hestia_login.DiscordClientSecret = tc.DiscordClientSecret
		hestia_login.SetConf(tc.DiscordClientID, tc.DiscordClientSecret)

	}
	log.Info().Msg("Hestia: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Hestia: PG connection connected")

	log.Info().Msg("Hestia: AWS Secrets Manager connection starting")
	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, awsAuthCfg)
	log.Info().Msg("Hestia: AWS Secrets Manager connected")

	inMemFsErr := artemis_validator_signature_service_routing.InitRouteMapInMemFS(ctx)
	if inMemFsErr != nil {
		log.Ctx(ctx).Err(inMemFsErr).Msg("Hydra: InitRouteMapInMemFS failed")
		misc.DelayedPanic(inMemFsErr)
	}
	//go func() {
	//	artemis_validator_signature_service_routing.InitAsyncServiceAuthRoutePollingHeartbeatAll(ctx)
	//}()

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

	log.Info().Msg("Hestia: InitHestiaQuickNodeWorker starting")
	quicknode_orchestrations.InitHestiaQuickNodeWorker(context.Background(), temporalAuthConfigHestia)
	cHestia := quicknode_orchestrations.HestiaQnWorker.Worker.ConnectTemporalClient()
	defer cHestia.Close()
	quicknode_orchestrations.HestiaQnWorker.Worker.RegisterWorker(cHestia)
	err = quicknode_orchestrations.HestiaQnWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Hestia: %s HestiaQnWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Hestia: InitHestiaQuickNodeWorker done")
	log.Info().Msg("Hestia: InitHestiaIrisPlatformServicesWorker start")
	platform_service_orchestrations.InitHestiaIrisPlatformServicesWorker(context.Background(), temporalAuthConfigHestia)
	platform_service_orchestrations.HestiaPlatformServiceWorker.Worker.RegisterWorker(cHestia)
	err = platform_service_orchestrations.HestiaPlatformServiceWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Hestia: %s HestiaPlatformServiceWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Hestia: InitHestiaIrisPlatformServicesWorker done")

	log.Info().Msg("Hestia: InitKronosWorker start")
	kronos_helix.InitKronosHelixWorker(context.Background(), temporalAuthConfigKronos)
	cKronos := kronos_helix.KronosServiceWorker.Worker.ConnectTemporalClient()
	defer cKronos.Close()
	kronos_helix.KronosServiceWorker.Worker.RegisterWorker(cKronos)
	err = kronos_helix.KronosServiceWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Hestia: %s InitKronosWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Hestia: InitKronosWorker done")

	if env == "local" || env == "production-local" {
		irisHost := "http://localhost:8080"
		srv.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"http://localhost:3000", "http://localhost:9002", irisHost, "http://promql.promql-edc89f30.svc.cluster.local",
				"https://accounts.google.com", "https://oauth2.googleapis.com", "http://localhost:3010",
				"https://twitter.com/", "https://twitter.com/i/oauth2/authorize",
				"https://api.twitter.com/2/oauth2", "https://api.twitter.com/2/oauth2/token",
				"https://hestia.zeus.fyi/social/v1/auth/twitter/callback", "https://cloud.zeus.fyi/social/v1/twitter/callback",
			},
			AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization,
				echo.HeaderAccessControlAllowHeaders, "X-CSRF-Token", "Accept-Encoding",
				SelectedRouteHeader, SelectedLatencyHeader, SelectedRouteGroupHeader, SelectedResponseReceivedAt, AdaptiveMetricsKey,
			},
			AllowCredentials: true,
		}))
		hestia_login.Domain = "localhost"
		hestia_billing.IrisApiUrl = irisHost
		v1_ethereum_aws.LambdaBaseDirIn = "/"
	} else {
		srv.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"https://cloud.zeus.fyi", "https://api.zeus.fyi", "https://hestia.zeus.fyi",
				"https://hestia.zeus.fyi/social/v1/auth/twitter/callback", "https://hestia.zeus.fyi/social/v1/twitter/callback",
				"https://cloud.zeus.fyi/social/v1/auth/twitter/callback", "https://cloud.zeus.fyi/social/v1/twitter/callback",
				"https://api.twitter.com/2/oauth2/token",
				"https://twitter.com/", "https://twitter.com/i/oauth2", "https://twitter.com/i/oauth2/authorize",
				"https://iris.zeus.fyi", "https://flows.zeus.fyi", "https://api.flows.zeus.fyi", "https://staging.api.flows.zeus.fyi",
				"https://quicknode.com", "https://staging.flows.zeus.fyi",
				"https://accounts.google.com", "https://oauth2.googleapis.com",
				"http://promql.promql-edc89f30.svc.cluster.local",
				"https://oauth.reddit.com",
				"https://hestia.zeus.fyi/reddit/callback",
				"https://cloud.zeus.fyi/reddit/callback",
			},
			AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization,
				echo.HeaderAccessControlAllowHeaders, "X-CSRF-Token", "Accept-Encoding",
				SelectedRouteHeader, SelectedLatencyHeader, SelectedRouteGroupHeader, SelectedResponseReceivedAt, AdaptiveMetricsKey,
			},
			AllowCredentials: true,
		}))
		hestia_login.Domain = "zeus.fyi"
		v1_ethereum_aws.LambdaBaseDirIn = "/etc/serverless/"
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
