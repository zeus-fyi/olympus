package tyche_server

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	v1_tyche "github.com/zeus-fyi/olympus/tyche/api/v1"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
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
	log.Info().Msgf("Tyche: Environment %s", env)
	SetConfigByEnv(ctx, env)
	srv.E = v1_tyche.Routes(srv.E)
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9000", "server port")
	Cmd.Flags().StringVar(&env, "env", "production-local", "environment")
	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "MEV Active Trading Module",
	Short: "Tyche",
	Run: func(cmd *cobra.Command, args []string) {
		Tyche()
	},
}
