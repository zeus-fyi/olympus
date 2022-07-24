package server

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/autok8s/api/v1"
)

var cfg = Config{}

func AutoK8s() {
	srv := NewAutoK8sServer(cfg)
	// Echo instance
	if cfg.K8sUtil.CfgPath == "" {
		log.Debug().Msg("AutoK8sCmd")
		log.Debug().Msg("The k8s config path was empty, so using default path")
		cfg.K8sUtil.CfgPath = cfg.K8sUtil.DefaultK8sCfgPath()
	}
	log.Debug().Msgf("The k8s config path %s:", cfg.K8sUtil.CfgPath)
	srv.E = v1.InitRouter(srv.E, cfg.K8sUtil)

	// Middleware
	srv.E.Use(middleware.Logger())
	srv.E.Use(middleware.Recover())
	// Start server
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9001", "server port")
	Cmd.Flags().StringVar(&cfg.K8sUtil.CfgPath, "kubie-config-path", "", "kubie config path")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "auto_k8s",
	Short: "A Transformer for K8s Actions",
	Run: func(cmd *cobra.Command, args []string) {
		AutoK8s()
	},
}
