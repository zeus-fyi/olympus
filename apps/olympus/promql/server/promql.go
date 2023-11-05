package promql_server

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1_promql "github.com/zeus-fyi/olympus/promql/api/v1"
)

var (
	cfg = Config{}
	env string
)

func PromQL() {
	ctx := context.Background()
	cfg.Host = "0.0.0.0"
	srv := NewPromQLServer(cfg)
	// Echo instance
	srv.E = v1_promql.Routes(srv.E)
	// Start server
	SetConfigByEnv(ctx, env)
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")
	Cmd.Flags().StringVar(&cfg.Port, "port", "8000", "server port")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Proxying prometheus",
	Short: "Proxying prometheus",
	Run: func(cmd *cobra.Command, args []string) {
		PromQL()
	},
}
