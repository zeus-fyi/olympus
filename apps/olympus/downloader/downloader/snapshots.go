package snapshot_init

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	init_jwt "github.com/zeus-fyi/zeus/pkg/aegis/jwt"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
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
	CloudCtxNsID      int    // id for infra location
	DataDir           filepaths.Path
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
	Cmd.Flags().IntVar(&Workload.CloudCtxNsID, "cloudCtxNsID", 0, "cloud ctx ns location info")
	Cmd.Flags().IntVar(&Workload.ProtocolNetworkID, "protocolNetworkID", 0, "identifier for protocol and network")
	Cmd.Flags().StringVar(&Workload.ClientName, "clientName", "", "client name")
	Cmd.Flags().StringVar(&Workload.WorkloadType, "workloadType", "", "workloadType") // eg validatorClient
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Downloads and extracts blockchain data and configs to your dataDir",
	Short: "Blockchain node data, and validator info, download procedure",
	Run: func(cmd *cobra.Command, args []string) {
		StartUp()
	},
}
