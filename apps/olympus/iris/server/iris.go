package iris_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1_iris "github.com/zeus-fyi/olympus/iris/api/v1"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	artemis_api_requests "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/api_requests"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

var (
	cfg             = Config{}
	temporalAuthCfg temporal_auth.TemporalAuth
	env             string
	authKeysCfg     auth_keys_config.AuthKeysCfg
)

func Iris() {
	cfg.Host = "0.0.0.0"
	// Echo instance
	ctx := context.Background()
	SetConfigByEnv(ctx, env)

	srv := NewIrisServer(cfg)
	srv.E = v1_iris.Routes(srv.E)
	// Start server
	log.Info().Msg("Iris: Starting ArtemisProxyWorker")
	c := artemis_api_requests.ArtemisProxyWorker.ConnectTemporalClient()
	defer c.Close()
	artemis_api_requests.ArtemisProxyWorker.Worker.RegisterWorker(c)
	err := artemis_api_requests.ArtemisProxyWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Iris: %s ArtemisProxyWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Iris: ArtemisProxyWorker Started")
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "8080", "server port")
	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")
	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Message proxy and router",
	Short: "Message proxy and router",
	Run: func(cmd *cobra.Command, args []string) {
		Iris()
	},
}
