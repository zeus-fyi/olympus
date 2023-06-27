package tyche_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	apollo_metrics_workload_info "github.com/zeus-fyi/olympus/pkg/apollo/metrics/workload_info"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	v1_tyche "github.com/zeus-fyi/olympus/tyche/api/v1"
	tyche_metrics "github.com/zeus-fyi/olympus/tyche/metrics"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

var (
	cfg             = Config{}
	temporalAuthCfg temporal_auth.TemporalAuth
	authKeysCfg     auth_keys_config.AuthKeysCfg
	env             string
	Workload        WorkloadInfo
	awsRegion       = "us-west-1"
)

type WorkloadInfo struct {
	zeus_common_types.CloudCtxNs
	ProtocolNetworkID int // eg. mainnet
}

func Tyche() {
	ctx := context.Background()
	cfg.Host = "0.0.0.0"
	srv := NewTycheServer(cfg)
	log.Ctx(ctx).Info().Msg("Tyche: Initializing configs by environment type")
	SetConfigByEnv(ctx, env)
	mwi := apollo_metrics_workload_info.NewWorkloadInfo("tyche", Workload.CloudCtxNs)
	tyche_metrics.InitTycheMetrics(ctx, mwi)
	srv.E = v1_tyche.Routes(srv.E)
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9000", "server port")
	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")
	Cmd.Flags().IntVar(&Workload.ProtocolNetworkID, "protocol-network-id", hestia_req_types.EthereumMainnetProtocolNetworkID, "identifier for protocol and network")

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
		Tyche()
	},
}
