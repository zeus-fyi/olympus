package hypnos_server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1_hypnos "github.com/zeus-fyi/olympus/hypnos/api/v1"
)

var cfg = Config{}

func Hypnos() {
	cfg.Host = "0.0.0.0"
	srv := NewHypnosServer(cfg)
	// Echo instance
	srv.E = v1_hypnos.Routes(srv.E)
	// Start server
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "8888", "server port")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Proxying a container, pod, network, or otherwise",
	Short: "The god of hypnosis",
	Run: func(cmd *cobra.Command, args []string) {
		Hypnos()
	},
}
