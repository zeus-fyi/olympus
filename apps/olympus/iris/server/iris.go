package iris_server

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	iris_api "github.com/zeus-fyi/olympus/iris/api"
	iris_metrics "github.com/zeus-fyi/olympus/iris/api/metrics"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	iris_api_requests "github.com/zeus-fyi/olympus/pkg/iris/proxy/orchestrations/api_requests"
	iris_serverless "github.com/zeus-fyi/olympus/pkg/iris/serverless"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

var (
	cfg             = Config{}
	temporalAuthCfg temporal_auth.TemporalAuth
	env             string
	authKeysCfg     auth_keys_config.AuthKeysCfg
)

const (
	SelectedRouteHeader        = "X-Selected-Route"
	SelectedLatencyHeader      = "X-Response-Latency-Milliseconds"
	SelectedRouteGroupHeader   = "X-Route-Group"
	SelectedResponseReceivedAt = "X-Response-Received-At-UTC"
	AdaptiveMetricsKey         = "X-Adaptive-Metrics-Key"
)

func Iris() {
	cfg.Host = "0.0.0.0"
	// Echo instance
	ctx := context.Background()
	SetConfigByEnv(ctx, env)

	srv := NewIrisServer(cfg)
	srv.E = iris_api.Routes(srv.E)

	metricsSrv := NewMetricsServer(cfg)
	metricsSrv.E = iris_api.MetricRoutes(metricsSrv.E)
	iris_metrics.InitIrisMetrics()
	// Start server
	log.Info().Msg("Iris: Starting IrisProxyWorker")
	c := iris_api_requests.IrisProxyWorker.ConnectTemporalClient()
	defer c.Close()
	iris_api_requests.IrisProxyWorker.Worker.RegisterWorker(c)
	err := iris_api_requests.IrisProxyWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Iris: %s IrisProxyWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Iris: Started IrisProxyWorker")
	log.Info().Msg("Iris: Starting IrisCacheWorker")
	c1 := iris_api_requests.IrisCacheWorker.ConnectTemporalClient()
	defer c1.Close()
	iris_api_requests.IrisCacheWorker.Worker.RegisterWorker(c1)
	err = iris_api_requests.IrisCacheWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Iris: %s IrisCacheWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Iris: IrisCacheWorker Started")

	log.Info().Msg("Iris: Starting InitIrisPlatformServicesWorker")
	c2 := iris_serverless.IrisPlatformServicesWorker.ConnectTemporalClient()
	defer c2.Close()
	iris_serverless.IrisPlatformServicesWorker.Worker.RegisterWorker(c2)
	err = iris_serverless.IrisPlatformServicesWorker.Worker.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Iris: %s IrisPlatformServicesWorker.Worker.Start failed", env)
		misc.DelayedPanic(err)
	}
	log.Info().Msg("Iris: Started InitIrisPlatformServicesWorker")

	go func() {
		metricsSrv.Start()
	}()

	if env == "local" || env == "production-local" {
		hestiaHost := "http://localhost:9002"
		srv.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"http://localhost:3000", hestiaHost, "http://localhost:8080"},
			AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderAccessControlAllowHeaders,
				"X-CSRF-Token", "Accept-Encoding", "X-Route-Group",
				iris_programmable_proxy_v1_beta.DurableExecutionID, iris_programmable_proxy_v1_beta.LoadBalancingStrategy,
				SelectedRouteHeader, SelectedLatencyHeader, SelectedRouteGroupHeader, SelectedResponseReceivedAt, AdaptiveMetricsKey,
			},
			AllowCredentials: true,
		}))
	} else {

		srv.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"https://cloud.zeus.fyi", "https://api.flows.zeus.fyi",
				"https://staging.api.flows.zeus.fyi",
				"https://api.zeus.fyi", "https://hestia.zeus.fyi", "https://iris.zeus.fyi", "https://quicknode.com",
				"https://oauth.reddit.com", "https://flows.zeus.fyi"},
			AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization,
				echo.HeaderAccessControlAllowHeaders, "X-CSRF-Token", "Accept-Encoding", "X-Route-Group",
				iris_programmable_proxy_v1_beta.DurableExecutionID, iris_programmable_proxy_v1_beta.LoadBalancingStrategy,
				SelectedRouteHeader, SelectedLatencyHeader, SelectedRouteGroupHeader, SelectedResponseReceivedAt, AdaptiveMetricsKey,
			},
			AllowCredentials: true,
		}))
	}
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
