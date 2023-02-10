package hydra_server

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1_hydra "github.com/zeus-fyi/olympus/hydra/api/v1"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

var (
	cfg             = Config{}
	temporalAuthCfg temporal_auth.TemporalAuth
	authKeysCfg     auth_keys_config.AuthKeysCfg
	env             string
	Workload        WorkloadInfo
)

// TODO set up some polling mechanism to update validators to send serverless functions to

type WorkloadInfo struct {
	zeus_common_types.CloudCtxNs
	ProtocolNetworkID int // eg. mainnet
}

func Hydra() {
	ctx := context.Background()
	cfg.Host = "0.0.0.0"
	srv := NewHydraServer(cfg)
	SetConfigByEnv(ctx, env)
	srv.E = v1_hydra.Routes(srv.E)
	// Start server
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9000", "server port")
	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")
	Cmd.Flags().IntVar(&Workload.ProtocolNetworkID, "protocol-network-id", 0, "identifier for protocol and network")

	Cmd.Flags().StringVar(&Workload.CloudCtxNs.CloudProvider, "cloud-provider", "", "cloud-provider")
	Cmd.Flags().StringVar(&Workload.CloudCtxNs.Context, "ctx", "", "context")
	Cmd.Flags().StringVar(&Workload.CloudCtxNs.Namespace, "ns", "", "namespace")
	Cmd.Flags().StringVar(&Workload.CloudCtxNs.Region, "region", "", "region")

	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Web3signing proxy router",
	Short: "Proxy",
	Run: func(cmd *cobra.Command, args []string) {
		Hydra()
	},
}
