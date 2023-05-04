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
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	topology_worker "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workers/topology"
	router "github.com/zeus-fyi/olympus/zeus/api"
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
		_, sw := auth_startup.RunZeusDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cmd := exec.Command("doctl", "auth", "init", "-t", sw.DoctlToken)
		err := cmd.Run()
		if err != nil {
			log.Fatal().Msg("RunDigitalOceanS3BucketObjSecretsProcedure: failed to auth doctl, shutting down the server")
			misc.DelayedPanic(err)
		}
		api_auth_temporal.InitOrchestrationDigitalOceanClient(ctx, sw.DoctlToken)
		api_auth_temporal.InitOrchestrationGcpClient(ctx, sw.GcpAuthJsonBytes)
		hestia_stripe.InitStripe(sw.StripeSecretKey)
	case "production-local":
		log.Info().Msg("Zeus: production local, auth procedure starting")
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.ProdLocalAuthKeysCfg)
		inMemFs := auth_startup.RunDigitalOceanS3BucketObjAuthProcedure(ctx, authCfg)
		cfg.K8sUtil.ConnectToK8sFromInMemFsCfgPath(inMemFs)
		temporalAuthCfg = tc.DevTemporalAuth
		_, sw := auth_startup.RunZeusDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		api_auth_temporal.InitOrchestrationDigitalOceanClient(ctx, sw.DoctlToken)
		api_auth_temporal.InitOrchestrationGcpClient(ctx, sw.GcpAuthJsonBytes)
		hestia_stripe.InitStripe(tc.StripeTestSecretAPIKey)
	case "local":
		log.Info().Msg("Zeus: local, auth procedure starting")
		tc := configs.InitLocalTestConfigs()
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.DevAuthKeysCfg)
		inMemFs := auth_startup.RunDigitalOceanS3BucketObjAuthProcedure(ctx, authCfg)
		cfg.K8sUtil.ConnectToK8sFromInMemFsCfgPath(inMemFs)
		temporalAuthCfg = tc.DevTemporalAuth
		api_auth_temporal.InitOrchestrationDigitalOceanClient(ctx, tc.DigitalOceanAPIKey)
		api_auth_temporal.InitOrchestrationGcpClient(ctx, tc.GcpAuthJson)
		hestia_stripe.InitStripe(tc.StripeTestSecretAPIKey)
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
