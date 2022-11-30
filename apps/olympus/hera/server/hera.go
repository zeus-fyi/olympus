package hera_server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1_hera "github.com/zeus-fyi/olympus/hera/api/v1"
)

var cfg = Config{}

func Hera() {
	cfg.Host = "0.0.0.0"
	srv := NewHeraServer(cfg)
	// Echo instance
	srv.E = v1_hera.Routes(srv.E)
	// Start server
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9008", "server port")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Generation of keys, networks, code, and more.",
	Short: "Generation of keys, networks, code..",
	Run: func(cmd *cobra.Command, args []string) {
		Hera()
	},
}
