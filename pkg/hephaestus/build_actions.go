package hephaestus_build_actions

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	athena_workloads "github.com/zeus-fyi/olympus/pkg/athena/workloads"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

var (
	appName     string
	dataDir     filepaths.Path
	authKeysCfg auth_keys_config.AuthKeysCfg
	env         string
	Workload    athena_workloads.WorkloadInfo
)

func StartUp() {
	dataDir = filepaths.Path{
		DirIn:  "/data",
		DirOut: "/data",
	}
	ctx := context.Background()
	cfg := Config{}
	switch env {
	case "production":
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunAthenaDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		dataDir.DirOut = "../"
	case "local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.DevAuthKeysCfg
		cfg.PGConnStr = tc.LocalDbPgconn
		dataDir.DirOut = "../"
	}
	Rebuild()
	Upload(ctx)
}

func Rebuild() {
	// should be able to do git pull, go build, and restart the app
	// for now, assume gitpull on init-container
	// so just rebuild the app and push up the binary

	sanitizedAppName := ""
	switch appName {
	case "iris":
		sanitizedAppName = appName
	case "hardhat":
		sanitizedAppName = appName
	case "artemis":
		sanitizedAppName = appName
	case "zeus":
		sanitizedAppName = appName
	default:
		log.Fatal().Msg("invalid app name")
	}
	dataDir.FnIn = fmt.Sprintf("apps/olympus/%s", sanitizedAppName)
	log.Info().Interface("appName", sanitizedAppName).Msg("Hephaestus Rebuild")
	rebuildCmd := "go"
	buildArgs := []string{
		"build",
		fmt.Sprintf("-ldflags=-s -w"),
		"-o",
		sanitizedAppName,
		dataDir.FileInPath(),
	}
	cmd := exec.Command(rebuildCmd, buildArgs...)
	err := cmd.Run()
	if err != nil {
		log.Fatal().Msg("failed to rebuild app")
		misc.DelayedPanic(err)
	}
	log.Info().Interface("appName", sanitizedAppName).Msg("Hephaestus Done Rebuild")
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&appName, "app", "", "app name to rebuild")
	Cmd.Flags().StringVar(&dataDir.DirOut, "dataDir", "/data", "data directory location")
	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
	Cmd.Flags().StringVar(&env, "env", "local", "environment")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "File syncs and fast app rebuilds",
	Short: "File syncs and app rebuilds",
	Run: func(cmd *cobra.Command, args []string) {
		StartUp()
	},
}
