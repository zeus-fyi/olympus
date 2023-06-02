package iris_server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1_iris "github.com/zeus-fyi/olympus/iris/api/v1"
)

var cfg = Config{}

func Iris() {
	cfg.Host = "0.0.0.0"
	srv := NewIrisServer(cfg)
	// Echo instance
	srv.E = v1_iris.Routes(srv.E)
	// Start server
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "8080", "server port")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Message proxy and router",
	Short: "Message proxy and router",
	Run: func(cmd *cobra.Command, args []string) {
		Iris()
	},
}
