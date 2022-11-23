package zeus_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	topology_worker "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workers/topology"
	router "github.com/zeus-fyi/olympus/zeus/api"
)

var cfg = Config{}
var authKeysCfg auth_keys_config.AuthKeysCfg
var temporalAuthCfg temporal_auth.TemporalAuth
var env string

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
		_, _ = auth_startup.RunDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
	case "production-local":
		log.Info().Msg("Zeus: production local, auth procedure starting")
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.ProdLocalAuthKeysCfg)
		inMemFs := auth_startup.RunDigitalOceanS3BucketObjAuthProcedure(ctx, authCfg)
		cfg.K8sUtil.ConnectToK8sFromInMemFsCfgPath(inMemFs)
		temporalAuthCfg = tc.ProdLocalTemporalAuth

		_, _ = auth_startup.RunDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
	case "local":
		log.Info().Msg("Zeus: local, auth procedure starting")
		tc := configs.InitLocalTestConfigs()
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.DevAuthKeysCfg)
		inMemFs := auth_startup.RunDigitalOceanS3BucketObjAuthProcedure(ctx, authCfg)
		cfg.K8sUtil.ConnectToK8sFromInMemFsCfgPath(inMemFs)
		temporalAuthCfg = tc.ProdLocalTemporalAuth
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
	srv.E = router.InitRouter(srv.E, cfg.K8sUtil)
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

	Cmd.Flags().StringVar(&env, "env", "local", "environment")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "zeus",
	Short: "A transformer for distributed infra actions",
	Run: func(cmd *cobra.Command, args []string) {
		Zeus()
	},
}
