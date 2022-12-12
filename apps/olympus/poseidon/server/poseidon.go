package poseidon_server

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
	"github.com/zeus-fyi/olympus/pkg/poseidon/poseidon_orchestrations"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	v1_poseidon "github.com/zeus-fyi/olympus/poseidon/api/v1"
)

var cfg = Config{}
var authKeysCfg auth_keys_config.AuthKeysCfg
var temporalAuthCfg temporal_auth.TemporalAuth
var env, bearer string

func Poseidon() {
	cfg.Host = "0.0.0.0"
	ctx := context.Background()
	srv := NewPoseidonServer(cfg)
	// Echo instance
	srv.E = v1_poseidon.Routes(srv.E)

	switch env {
	case "production":
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunPoseidonDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
		temporalAuthCfg = temporal_auth.TemporalAuth{
			ClientCertPath:   "/etc/ssl/certs/ca.pem",
			ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
			Namespace:        "production-poseidon.ngb72",
			HostPort:         "production-poseidon.ngb72.tmprl.cloud:7233",
		}
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		temporalAuthCfg = tc.ProdLocalTemporalAuthPoseidon
	case "local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.DevAuthKeysCfg
		cfg.PGConnStr = tc.LocalDbPgconn
		temporalAuthCfg = tc.ProdLocalTemporalAuthPoseidon
	}

	log.Info().Msg("Poseidon: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Poseidon: PG connected")
	bearer = auth_startup.FetchTemporalAuthBearer(ctx)

	log.Info().Msgf("Poseidon: %s temporal auth and init procedure starting", env)
	poseidon_orchestrations.PoseidonBearer = bearer
	poseidon_orchestrations.InitPoseidonWorker(ctx, temporalAuthCfg)
	c := poseidon_orchestrations.PoseidonSyncWorker.TemporalClient.ConnectTemporalClient()
	defer c.Close()
	poseidon_orchestrations.PoseidonSyncWorker.Worker.RegisterWorker(c)
	err := poseidon_orchestrations.PoseidonSyncWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Poseidon: %s topology_worker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	poseidon_orchestrations.PoseidonS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg)

	log.Info().Msgf("Poseidon: %s temporal setup is complete", env)
	log.Info().Msgf("Poseidon: %s server starting", env)
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9006", "server port")
	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")
	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Downloading service for blockchain data",
	Short: "Chain download service",
	Run: func(cmd *cobra.Command, args []string) {
		Poseidon()
	},
}
