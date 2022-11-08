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
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	topology_worker "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workers/topology"
	router "github.com/zeus-fyi/olympus/zeus/api"
)

var cfg = Config{}
var authKeysCfg auth_keys_config.AuthKeysCfg
var env string

func Zeus() {
	log.Info().Msg("Zeus: starting")
	srv := NewZeusServer(cfg)
	// Echo instance
	ctx := context.Background()
	switch env {
	case "production":
		log.Info().Msg("Zeus: production auth procedure starting")
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		inMemFs := auth_startup.RunDigitalOceanS3BucketObjAuthProcedure(ctx, authCfg)
		cfg.K8sUtil.ConnectToK8sFromInMemFsCfgPath(inMemFs)

		temporalAuthCfg := temporal_base.TemporalAuth{
			ClientCertPath:   "/etc/ssl/certs/ca.pem",
			ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
			Namespace:        "production-zeus.ngb72",
			HostPort:         "production-zeus.ngb72.tmprl.cloud:7233",
		}
		_, _ = topology_worker.InitTopologyWorker(temporalAuthCfg)
	case "production-local":
		log.Info().Msg("Zeus: production local, auth procedure starting")
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.ProdLocalAuthKeysCfg)
		inMemFs := auth_startup.RunDigitalOceanS3BucketObjAuthProcedure(ctx, authCfg)
		cfg.K8sUtil.ConnectToK8sFromInMemFsCfgPath(inMemFs)
		_, _ = topology_worker.InitTopologyWorker(tc.ProdLocalTemporalAuth)
	case "local":
		tc := configs.InitLocalTestConfigs()
		authCfg := auth_startup.NewDefaultAuthClient(ctx, tc.DevAuthKeysCfg)
		inMemFs := auth_startup.RunDigitalOceanS3BucketObjAuthProcedure(ctx, authCfg)
		cfg.K8sUtil.ConnectToK8sFromInMemFsCfgPath(inMemFs)
		_, _ = topology_worker.InitTopologyWorker(tc.DevTemporalAuth)
	}

	log.Info().Msg("Zeus: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)

	srv.E = router.InitRouter(srv.E, cfg.K8sUtil)

	log.Info().Msg("Zeus: production server starting")
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
	Short: "A transformer for distributed infra actions",
	Run: func(cmd *cobra.Command, args []string) {
		Zeus()
	},
}
