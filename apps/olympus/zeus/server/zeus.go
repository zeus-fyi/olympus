package server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	router "github.com/zeus-fyi/olympus/zeus/api"
)

var cfg = Config{}

func Zeus() {
	srv := NewZeusServer(cfg)
	// Echo instance
	if cfg.K8sUtil.CfgPath == "" {
		log.Debug().Msg("ZeusCmd")
		log.Debug().Msg("The k8s config path was empty, so using default path")
		cfg.K8sUtil.CfgPath = cfg.K8sUtil.DefaultK8sCfgPath()
	}
	log.Debug().Msgf("The k8s config path %s:", cfg.K8sUtil.CfgPath)
	srv.E = router.InitRouter(srv.E, cfg.K8sUtil)

	ctx := context.Background()
	apps.Pg = apps.Db{}
	apps.Pg.InitPG(ctx, cfg.PGConnStr)

	// Start server
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9001", "server port")
	Cmd.Flags().StringVar(&cfg.K8sUtil.CfgPath, "kubie-config-path", "", "kubie config path")
	Cmd.Flags().StringVar(&cfg.PGConnStr, "postgres-conn-str", "", "postgres connection string")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "zeus",
	Short: "A transformer for distributed infra actions",
	Run: func(cmd *cobra.Command, args []string) {
		Zeus()
	},
}
