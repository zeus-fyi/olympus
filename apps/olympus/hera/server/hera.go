package hera_server

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	v1_hera "github.com/zeus-fyi/olympus/hera/api/v1"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
)

var cfg = Config{}
var authKeysCfg auth_keys_config.AuthKeysCfg
var temporalAuthCfg temporal_auth.TemporalAuth
var env string

func Hera() {
	ctx := context.Background()

	cfg.Host = "0.0.0.0"
	srv := NewHeraServer(cfg)
	// Echo instance
	srv.E = v1_hera.Routes(srv.E)

	switch env {
	case "production":
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunHeraDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
		temporalAuthCfg = temporal_auth.TemporalAuth{
			ClientCertPath:   "/etc/ssl/certs/ca.pem",
			ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
			Namespace:        "production-hera.ngb72",
			HostPort:         "production-hera.ngb72.tmprl.cloud:7233",
		}
	case "production-local":
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunHeraDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
		hera_openai.InitHeraOpenAI(sw.OpenAIToken)
	case "local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.DevAuthKeysCfg
		cfg.PGConnStr = tc.ProdLocalDbPgconn
		hera_openai.InitHeraOpenAI(tc.OpenAIAuth)
	}

	if env == "local" || env == "production-local" {
		srv.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderAccessControlAllowHeaders, "X-CSRF-Token", "Accept-Encoding"},
			AllowCredentials: true,
		}))
	} else {
		srv.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"https://cloud.zeus.fyi", "https://api.zeus.fyi", "https://hestia.zeus.fyi", "https://api.flows.zeus.fyi", "https://staging.flows.api.zeus.fyi"},
			AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderAccessControlAllowHeaders, "X-CSRF-Token", "Accept-Encoding"},
			AllowCredentials: true,
		}))
	}
	log.Info().Msg("Hera: PG connection starting")
	apps.Pg.InitPG(ctx, cfg.PGConnStr)
	srv.Start()
}

func init() {
	viper.AutomaticEnv()
	Cmd.Flags().StringVar(&cfg.Port, "port", "9008", "server port")
	Cmd.Flags().StringVar(&env, "env", "local", "environment")
	Cmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	Cmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	Cmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
}

// Cmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use:   "Generation of keys, networks, code, and more.",
	Short: "Generation of keys, networks, code..",
	Run: func(cmd *cobra.Command, args []string) {
		Hera()
	},
}
