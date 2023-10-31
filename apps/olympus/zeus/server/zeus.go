package zeus_server

import (
	"context"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/dynamic_secrets"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	topology_worker "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workers/topology"
	pods_workflows "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/pods"
	router "github.com/zeus-fyi/olympus/zeus/api"
	read_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/read"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

var (
	cfg             = Config{}
	authKeysCfg     auth_keys_config.AuthKeysCfg
	temporalAuthCfg temporal_auth.TemporalAuth
	env             string
)

func Zeus() {
	log.Info().Msg("Zeus: starting")
	cfg.Host = "0.0.0.0"
	srv := NewZeusServer(cfg)
	// Echo instance
	ctx := context.Background()
	switch env {
	case "production":
		log.Info().Msg("Zeus: production auth procedure starting")
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		inMemFs := auth_startup.RunDigitalOceanS3BucketObjAuthProcedure(ctx, authCfg)
		cfg.K8sUtil.ConnectToK8sFromInMemFsCfgPath(inMemFs)

		temporalAuthCfg = temporal_auth.TemporalAuth{
			ClientCertPath:   "/etc/ssl/certs/ca.pem",
			ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
			Namespace:        "production-zeus.ngb72",
			HostPort:         "production-zeus.ngb72.tmprl.cloud:7233",
		}
		dynMemFs, sw := auth_startup.RunZeusDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		dynamic_secrets.AegisInMemSecrets = dynMemFs
		cmd := exec.Command("doctl", "auth", "init", "-t", sw.DoctlToken)
		err := cmd.Run()
		if err != nil {
			log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: failed to auth doctl, shutting down the server")
			misc.DelayedPanic(err)
		}
		p := filepaths.Path{
			PackageName: "",
			DirIn:       "/secrets",
			DirOut:      "/secrets",
			FnOut:       "gcp_auth.json",
			FnIn:        "gcp_auth.json",
			Env:         "",
			FilterFiles: nil,
		}
		err = p.WriteToFileOutPath(sw.GcpAuthJsonBytes)
		if err != nil {
			log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: failed to write gcp auth json, shutting down the server")
			misc.DelayedPanic(err)
		}
		cmd = exec.Command("/google-cloud-sdk/bin/gcloud", "auth", "login", "--cred-file", "/secrets/gcp_auth.json")
		err = cmd.Run()
		if err != nil {
			log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: failed to auth gcloud, shutting down the server")
			misc.DelayedPanic(err)
		}
		err = p.RemoveFileInPath()
		if err != nil {
			log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: failed to remove gcp auth json, shutting down the server")
			misc.DelayedPanic(err)
		}
		api_auth_temporal.InitOrchestrationOvhCloudClient(ctx, sw.OvhAppKey, sw.OvhSecretKey, sw.OvhConsumerKey)
		api_auth_temporal.InitOrchestrationDigitalOceanClient(ctx, sw.DoctlToken)
		api_auth_temporal.InitOrchestrationGcpClient(ctx, sw.GcpAuthJsonBytes)
		api_auth_temporal.InitOrchestrationEksClient(ctx, sw.EksAuthAWS)
		hestia_stripe.InitStripe(sw.StripeSecretKey)
		cfg.PGConnStr = sw.PostgresAuth
	case "production-local":
		log.Info().Msg("Zeus: production local, auth procedure starting")
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.ProdLocalAuthKeysCfg)
		inMemFs := auth_startup.RunDigitalOceanS3BucketObjAuthProcedure(ctx, authCfg)
		cfg.K8sUtil.ConnectToK8sFromInMemFsCfgPath(inMemFs)
		temporalAuthCfg = tc.DevTemporalAuth
		dynMemFs, sw := auth_startup.RunZeusDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		dynamic_secrets.AegisInMemSecrets = dynMemFs
		api_auth_temporal.InitOrchestrationOvhCloudClient(ctx, sw.OvhAppKey, sw.OvhSecretKey, sw.OvhConsumerKey)
		api_auth_temporal.InitOrchestrationDigitalOceanClient(ctx, sw.DoctlToken)
		api_auth_temporal.InitOrchestrationGcpClient(ctx, sw.GcpAuthJsonBytes)
		api_auth_temporal.InitOrchestrationEksClient(ctx, sw.EksAuthAWS)
		hestia_stripe.InitStripe(tc.StripeTestSecretAPIKey)

		read_infra.CookbooksDirIn = "/Users/alex/go/Olympus/olympus/apps/zeus/cookbooks"
	case "local":
		log.Info().Msg("Zeus: local, auth procedure starting")
		tc := configs.InitLocalTestConfigs()
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.DevAuthKeysCfg)
		inMemFs := auth_startup.RunDigitalOceanS3BucketObjAuthProcedure(ctx, authCfg)
		cfg.K8sUtil.ConnectToK8sFromInMemFsCfgPath(inMemFs)
		temporalAuthCfg = tc.DevTemporalAuth
		api_auth_temporal.InitOrchestrationOvhCloudClient(ctx, tc.OvhAppKey, tc.OvhSecretKey, tc.OvhConsumerKey)
		api_auth_temporal.InitOrchestrationDigitalOceanClient(ctx, tc.DigitalOceanAPIKey)
		api_auth_temporal.InitOrchestrationGcpClient(ctx, tc.GcpAuthJson)
		hestia_stripe.InitStripe(tc.StripeTestSecretAPIKey)
		eksAuth := aegis_aws_auth.AuthAWS{
			AccountNumber: "",
			Region:        "us-west-1",
			AccessKey:     tc.AwsAccessKeyEks,
			SecretKey:     tc.AwsSecretKeyEks,
		}
		api_auth_temporal.InitOrchestrationEksClient(ctx, eksAuth)
		dynMemFs, _ := auth_startup.RunZeusDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		dynamic_secrets.AegisInMemSecrets = dynMemFs
	}

	log.Info().Msg("Zeus: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)

	log.Info().Msgf("Zeus: %s temporal auth and init procedure starting", env)
	api_auth_temporal.Bearer = auth_startup.FetchTemporalAuthBearer(ctx)
	topology_worker.InitTopologyWorker(temporalAuthCfg)

	c := topology_worker.Worker.TemporalClient.ConnectTemporalClient()
	defer c.Close()
	topology_worker.Worker.RegisterWorker(c)
	err := topology_worker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Zeus: %s topology_worker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}

	pods_workflows.InitPodsWorker(temporalAuthCfg)
	c2 := pods_workflows.PodsServiceWorker.TemporalClient.ConnectTemporalClient()
	defer c2.Close()
	pods_workflows.PodsServiceWorker.RegisterWorker(c2)
	err = pods_workflows.PodsServiceWorker.Start()

	log.Info().Msgf("Zeus: %s temporal setup is complete", env)
	log.Info().Msgf("Zeus: %s server starting", env)

	if env == "local" || env == "production-local" {
		mw := middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Access-Control-Allow-Headers", "X-CSRF-Token", "Accept-Encoding", "CloudCtxNsID"},
			AllowCredentials: true,
		})
		srv.E = router.InitRouter(srv.E, cfg.K8sUtil, mw)
	} else {
		mw := middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"https://cloud.zeus.fyi", "https://api.zeus.fyi", "https://hestia.zeus.fyi"},
			AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Access-Control-Allow-Headers", "X-CSRF-Token", "Accept-Encoding", "CloudCtxNsID"},
			AllowCredentials: true,
		})
		srv.E = router.InitRouter(srv.E, cfg.K8sUtil, mw)
	}
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9001", "server port")
	Cmd.Flags().StringVar(&cfg.PGConnStr, "postgres-conn-str", "postgresql://localhost/postgres?user=postgres&password=postgres", "postgres connection string")

	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")

	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "zeus",
	Short: "An orchestration engine for distributed infra actions",
	Run: func(cmd *cobra.Command, args []string) {
		Zeus()
	},
}
