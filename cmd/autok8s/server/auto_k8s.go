package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/autok8s/api/v1"
)

var cfg = Config{}

func AutoK8s() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Start server
	err := e.Start(":9000")
	if err != nil {
		log.Err(err)
	}
}

func init() {
	viper.AutomaticEnv()
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "auto_k8s",
	Short: "A Transformer for K8s Actions",
	Run: func(cmd *cobra.Command, args []string) {
		// now add routes (allows reuse of server code for tests by adding top router here)
		srv := NewAutoK8sServer(cfg)
		// Start server

		if cfg.K8sUtil.CfgPath == "" {
			log.Debug().Msg("AutoK8sCmd")
			log.Debug().Msg("The k8s config path was empty, so using default path")
			cfg.K8sUtil.ConnectToK8s()
		} else {
			log.Debug().Msgf("The k8s config path %s:", cfg.K8sUtil.CfgPath)
			cfg.K8sUtil.ConnectToK8sFromConfig(cfg.K8sUtil.CfgPath)
		}

		srv.E = v1.InitRouter(srv.E, cfg.K8sUtil)
		srv.Start()
	},
}
