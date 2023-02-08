package athena_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	athena_router "github.com/zeus-fyi/olympus/athena/api"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	"github.com/zeus-fyi/olympus/pkg/athena"
	athena_workloads "github.com/zeus-fyi/olympus/pkg/athena/workloads"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

var (
	cfg         = Config{}
	authKeysCfg auth_keys_config.AuthKeysCfg
	env         string
	dataDir     filepaths.Path
	Workload    athena_workloads.WorkloadInfo
)

func Athena() {
	ctx := context.Background()
	cfg.Host = "0.0.0.0"
	srv := NewAthenaServer(cfg)
	log.Info().Msgf("Athena: %s auth procedure starting", env)
	log.Info().Interface("workload", Workload).Msg("Athena: WorkloadInfo")
	switch env {
	case "production":
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunAthenaDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		dataDir.DirOut = "../"
	case "local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.DevAuthKeysCfg
		cfg.PGConnStr = tc.LocalDbPgconn
		dataDir.DirOut = "../"
	}
	log.Info().Msg("Athena: DigitalOceanS3AuthClient starting")
	athena.AthenaS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg)
	log.Info().Msg("Athena: DigitalOceanS3AuthClient done")

	log.Info().Msg("Athena: PG connection starting")
	apps.Pg = apps.Db{}
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	srv.E = athena_router.Routes(srv.E, dataDir)
	StartAndConfigClientNetworkSettings(ctx, Workload.ProtocolNetworkID, Workload.ClientName)
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9003", "server port")
	Cmd.Flags().StringVar(&dataDir.DirOut, "dataDir", "/data", "data directory location")
	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
	Cmd.Flags().StringVar(&env, "env", "local", "environment")

	Cmd.Flags().StringVar(&dataDir.DirIn, "dataDirIn", "/data", "data directory location")

	Cmd.Flags().StringVar(&Workload.WorkloadType, "workload-type", "", "workloadType") // eg validatorClient
	Cmd.Flags().IntVar(&Workload.ProtocolNetworkID, "protocol-network-id", 0, "identifier for protocol and network")
	Cmd.Flags().IntVar(&Workload.ReplicaCountNum, "replica-count-num", 0, "stateful set ordinal index")

	Cmd.Flags().StringVar(&Workload.CloudCtxNs.CloudProvider, "cloud-provider", "", "cloud-provider")
	Cmd.Flags().StringVar(&Workload.CloudCtxNs.Context, "ctx", "", "context")
	Cmd.Flags().StringVar(&Workload.CloudCtxNs.Namespace, "ns", "", "namespace")
	Cmd.Flags().StringVar(&Workload.CloudCtxNs.Region, "region", "", "region")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Web3 Middleware",
	Short: "A web3 infra middleware manager for apps on Olympus",
	Run: func(cmd *cobra.Command, args []string) {
		Athena()
	},
}
