package aegis_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1_aegis "github.com/zeus-fyi/olympus/aegis/api/v1"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

var (
	cfg                      = Config{}
	temporalAuthConfigKronos = temporal_auth.TemporalAuth{
		ClientCertPath:   "/etc/ssl/certs/ca.pem",
		ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
		Namespace:        "kronos.ngb72",
		HostPort:         "kronos.ngb72.tmprl.cloud:7233",
	}
	authKeysCfg auth_keys_config.AuthKeysCfg
	env         string
	dataDir     filepaths.Path
)

func Aegis() {
	cfg.Host = "0.0.0.0"
	srv := NewAegisServer(cfg)
	// Echo instance
	ctx := context.Background()
	switch env {
	case "production":
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
		hera_openai.InitHeraOpenAI(sw.OpenAIToken)
		kronos_helix.InitPagerDutyAlertClient(sw.PagerDutyApiKey)
		if (sw.PagerDutyApiKey == "") || (sw.PagerDutyRoutingKey == "") {
			panic("PAGERDUTY_API_KEY or PAGERDUTY_ROUTING_KEY is empty")
		}
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		kronos_helix.InitPagerDutyAlertClient(tc.PagerDutyApiKey)
		hera_openai.InitHeraOpenAI(tc.OpenAIAuth)
		temporalAuthConfigKronos = tc.DevTemporalAuth
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		dataDir.DirOut = "../"
	case "local":
		tc := configs.InitLocalTestConfigs()
		kronos_helix.InitPagerDutyAlertClient(tc.PagerDutyApiKey)
		hera_openai.InitHeraOpenAI(tc.OpenAIAuth)
		temporalAuthConfigKronos = tc.DevTemporalAuth
		authKeysCfg = tc.DevAuthKeysCfg
		cfg.PGConnStr = tc.LocalDbPgconn
		dataDir.DirOut = "../"
	}

	log.Info().Msg("Aegis: PG connection starting")
	apps.Pg = apps.Db{}
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	srv.E = v1_aegis.Routes(srv.E)

	log.Info().Msgf("Hestia %s artemis orchestration retrieving auth token", env)
	artemis_orchestration_auth.Bearer = auth_startup.FetchTemporalAuthBearer(ctx)
	log.Info().Msgf("Hestia %s artemis orchestration retrieving auth token done", env)

	log.Info().Msg("Aegis: InitKronosWorker start")
	kronos_helix.InitKronosHelixWorker(context.Background(), temporalAuthConfigKronos)
	cKronos := kronos_helix.KronosServiceWorker.Worker.ConnectTemporalClient()
	defer cKronos.Close()
	kronos_helix.KronosServiceWorker.Worker.RegisterWorker(cKronos)
	err := kronos_helix.KronosServiceWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Hestia: %s InitKronosWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Aegis: InitKronosWorker done")

	// Start server
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9007", "server port")
	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Authorizing requests to services",
	Short: "Authorization",
	Run: func(cmd *cobra.Command, args []string) {
		Aegis()
	},
}
