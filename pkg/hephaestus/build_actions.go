package hephaestus_build_actions

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appName string
)

func Rebuild() {
	// should be able to do git pull, go build, and restart the app
	// for now, assume gitpull on init-container
	// so just rebuild the app and push up the binary
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&appName, "app", "", "app name to rebuild")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "File syncs and fast app rebuilds",
	Short: "File syncs and app rebuilds",
	Run: func(cmd *cobra.Command, args []string) {
		Rebuild()
	},
}
