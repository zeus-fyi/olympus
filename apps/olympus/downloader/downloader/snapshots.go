package olympus_snapshot_init

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	"github.com/zeus-fyi/olympus/pkg/athena"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

var (
	env             string
	onlyIfEmptyDir  bool
	compressionType string
	jwtToken        string
	useDefaultToken bool
	Workload        WorkloadInfo
	authKeysCfg     auth_keys_config.AuthKeysCfg
	cfg             = Config{}
)

type Config struct {
	PGConnStr string
}

type WorkloadInfo struct {
	WorkloadType      string // eg, validatorClient, beaconExecClient, beaconConsensusClient
	ClientName        string // eg. lighthouse, geth
	ProtocolNetworkID int    // eg. mainnet
	ReplicaCountNum   int    // eg. stateful set ordinal index
	zeus_common_types.CloudCtxNs
	DataDir filepaths.Path
}

func StartUp() {
	ctx := context.Background()
	log.Ctx(ctx).Info().Interface("workload", Workload)

	log.Info().Msg("Downloader: DigitalOceanS3AuthClient starting")
	SetConfigByEnv(ctx, env)
	log.Info().Msg("Downloader: DigitalOceanS3AuthClient done")

	log.Info().Msg("Downloader: NewDigitalOceanS3AuthClient connecting")
	athena.AthenaS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg)
	log.Info().Msg("Downloader: NewDigitalOceanS3AuthClient done")

	log.Info().Msg("Downloader: InitPG connecting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Downloader: InitPG done")

	log.Info().Msg("Downloader: InitWorkloadAction starting")
	InitWorkloadAction(ctx, Workload)
	log.Info().Msg("Downloader: InitWorkloadAction done")
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().BoolVar(&onlyIfEmptyDir, "onlyIfEmptyDir", true, "only download & extract if the datadir is empty")
	Cmd.Flags().StringVar(&compressionType, "compressionExtension", ".tar.lz4", "compression type")
	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")

	// ethereum
	Cmd.Flags().StringVar(&jwtToken, "jwt", "0x6ad1acdc50a4141e518161ab2fe2bf6294de4b4d48bf3582f22cae8113f0cadc", "set jwt in datadir")
	Cmd.Flags().BoolVar(&useDefaultToken, "useDefaultToken", true, "use default jwt token")

	// internal
	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")

	// workload info
	Cmd.Flags().StringVar(&Workload.DataDir.DirIn, "dataDir", "/data", "data directory location")
	Cmd.Flags().StringVar(&Workload.WorkloadType, "workload-type", "", "workloadType") // eg validatorClient
	Cmd.Flags().StringVar(&Workload.ClientName, "client-name", "", "client name")
	Cmd.Flags().IntVar(&Workload.ReplicaCountNum, "replica-count-num", 0, "stateful set ordinal index")
	Cmd.Flags().IntVar(&Workload.ProtocolNetworkID, "protocol-network-id", 0, "identifier for protocol and network")

	Cmd.Flags().StringVar(&Workload.CloudCtxNs.CloudProvider, "cloud-provider", "", "cloud-provider")
	Cmd.Flags().StringVar(&Workload.CloudCtxNs.Context, "ctx", "", "context")
	Cmd.Flags().StringVar(&Workload.CloudCtxNs.Namespace, "ns", "", "namespace")
	Cmd.Flags().StringVar(&Workload.CloudCtxNs.Region, "region", "", "region")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Downloads and extracts blockchain data and configs to your dataDir",
	Short: "Blockchain node data, and validator info, download procedure",
	Run: func(cmd *cobra.Command, args []string) {
		StartUp()
	},
}
