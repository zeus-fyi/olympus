package snapshot_init

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
	preSignedURL    string
	env             string
	onlyIfEmptyDir  bool
	compressionType string
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
	log.Info().Msg("Downloader: DigitalOceanS3AuthClient starting")
	athena.AthenaS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg)
	SetConfigByEnv(ctx, env)
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	log.Info().Msg("Downloader: DigitalOceanS3AuthClient done")
	InitWorkloadAction(ctx, Workload)
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&preSignedURL, "downloadURL", "", "use a presigned bucket url")
	Cmd.Flags().BoolVar(&onlyIfEmptyDir, "onlyIfEmptyDir", true, "only download & extract if the datadir is empty")
	Cmd.Flags().StringVar(&compressionType, "compressionExtension", ".tar.lz4", "compression type")
	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")

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
