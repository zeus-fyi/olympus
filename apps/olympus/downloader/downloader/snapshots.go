package snapshot_init

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	init_jwt "github.com/zeus-fyi/zeus/pkg/aegis/jwt"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

var (
	preSignedURL    string
	onlyIfEmptyDir  bool
	compressionType string
	jwtToken        string
	useDefaultToken bool
	Workload        WorkloadInfo
)

type WorkloadInfo struct {
	WorkloadType      string // eg, validatorClient
	ClientName        string // eg. lighthouse, geth
	ProtocolNetworkID int    // eg. mainnet
	ReplicaCountNum   int    // eg. stateful set ordinal index
	zeus_common_types.CloudCtxNs
	DataDir filepaths.Path
}

func StartUp() {
	ctx := context.Background()
	InitWorkloadAction(ctx, Workload)
	if useDefaultToken {
		_ = init_jwt.SetTokenToDefault(Workload.DataDir, "jwt.hex", jwtToken)
	}
	ChainDownload()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&preSignedURL, "downloadURL", "", "use a presigned bucket url")
	Cmd.Flags().BoolVar(&onlyIfEmptyDir, "onlyIfEmptyDir", true, "only download & extract if the datadir is empty")
	Cmd.Flags().StringVar(&compressionType, "compressionExtension", ".tar.lz4", "compression type")
	Cmd.Flags().StringVar(&jwtToken, "jwt", "0x6ad1acdc50a4141e518161ab2fe2bf6294de4b4d48bf3582f22cae8113f0cadc", "set jwt in datadir")
	Cmd.Flags().BoolVar(&useDefaultToken, "useDefaultToken", true, "use default jwt token")

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
