package hephaestus_server

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	athena_workloads "github.com/zeus-fyi/olympus/pkg/athena/workloads"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

var (
	authKeysCfg auth_keys_config.AuthKeysCfg
	env         string
	dataDir     filepaths.Path
	Workload    athena_workloads.WorkloadInfo
	age         encryption.Age
	cfg         = Config{}
)

func createFormattedString(platform, appID, versionString, redditUsername string) string {
	return fmt.Sprintf("%s:%s:%s (by /u/%s)", platform, appID, versionString, redditUsername)
}
func Hephaestus() {
	ctx := context.Background()
	ro, err := reddit.NewReadonlyClient(reddit.WithUserAgent(createFormattedString("web", "hephaestus", "0.0.1", "zeus-fyi")))
	if err != nil {
		panic(err)
	}
	lpo := &reddit.ListOptions{Limit: 100}
	posts, resp, err := ro.Subreddit.NewPosts(ctx, Workload.WorkloadType, lpo)
	if err != nil {
		panic(err)
	}
	log.Info().Interface("posts", posts).Interface("resp", resp).Msg("Hephaestus")
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9000", "server port")
	Cmd.Flags().StringVar(&dataDir.DirOut, "dataDir", "/data", "data directory location")
	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
	Cmd.Flags().StringVar(&env, "env", "local", "environment")

	Cmd.Flags().StringVar(&dataDir.DirIn, "dataDirIn", "/data", "data directory location")

	Cmd.Flags().StringVar(&Workload.WorkloadType, "workload-type", "kubernetes", "workloadType") // eg validatorClient

}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Middleware",
	Short: "A infra middleware manager for apps on Olympus",
	Run: func(cmd *cobra.Command, args []string) {
		Hephaestus()
	},
}
