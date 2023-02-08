package hydra_choreography

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1_hydra_choreography "github.com/zeus-fyi/olympus/choreography/hydra/api/v1"
	hydra_choreography_metrics "github.com/zeus-fyi/olympus/choreography/hydra/metrics"
	apollo_metrics_workload_info "github.com/zeus-fyi/olympus/pkg/apollo/metrics/workload_info"
	zeus_client "github.com/zeus-fyi/zeus/pkg/zeus/client"
)

var (
	cfg    = Config{}
	bearer string
)

func Api() {
	ctx := context.Background()
	cfg.Host = "0.0.0.0"
	srv := NewChoreography(cfg)
	v1_hydra_choreography.ZeusClient = zeus_client.NewDefaultZeusClient(bearer)
	srv.E = v1_hydra_choreography.Routes(srv.E)

	// for beacon monitoring, todo use switch cases to build relevant metrics

	// todo add network name to metrics, workload type etc

	mwi := apollo_metrics_workload_info.NewWorkloadInfo("hydra", v1_hydra_choreography.CloudCtxNs)
	hydra_choreography_metrics.InitHydraMetrics(ctx, mwi)
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9999", "server port")
	// injected values
	Cmd.Flags().StringVar(&bearer, "bearer", "", "bearer")
	Cmd.Flags().StringVar(&v1_hydra_choreography.CloudCtxNs.CloudProvider, "cloud-provider", "", "cloud-provider")
	Cmd.Flags().StringVar(&v1_hydra_choreography.CloudCtxNs.Context, "ctx", "", "context")
	Cmd.Flags().StringVar(&v1_hydra_choreography.CloudCtxNs.Namespace, "ns", "", "namespace")
	Cmd.Flags().StringVar(&v1_hydra_choreography.CloudCtxNs.Region, "region", "", "region")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Hydra choreography actions",
	Short: "Hydra choreography",
	Run: func(cmd *cobra.Command, args []string) {
		Api()
	},
}
