package server

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "github.com/zeus-fyi/olympus/beacon-indexer/api/v1"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_fetcher"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/admin"
	redis_app "github.com/zeus-fyi/olympus/datastores/redis/apps"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	apollo_buckets "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/buckets"
)

var (
	RedisEndpointURL  string
	BeaconEndpointURL string
	PGConnStr         string
	env               string
	authKeysCfg       auth_keys_config.AuthKeysCfg
)

func Api() {
	//Echo instance
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	e := echo.New()
	e = v1.Routes(e)
	ctx := context.Background()

	switch env {
	case "production":
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunApolloDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		PGConnStr = sw.PostgresAuth
		BeaconEndpointURL = sw.MainnetBeaconURL
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		PGConnStr = tc.ProdLocalDbPgconn
		BeaconEndpointURL = tc.MainnetNodeUrl
	case "local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.DevAuthKeysCfg
		PGConnStr = tc.LocalDbPgconn
		BeaconEndpointURL = tc.MainnetNodeUrl
	}

	apps.Pg = apps.Db{}
	MaxConn := int32(10)
	MinConn := int32(3)
	MaxConnLifetime := 15 * time.Minute

	pgCfg := admin.ConfigChangePG{
		MaxConns:          &MaxConn,
		MinConn:           &MinConn,
		MaxConnLifetime:   &MaxConnLifetime,
		HealthCheckPeriod: nil,
	}
	apps.Pg.InitPG(ctx, PGConnStr)
	_ = admin.UpdateConfigPG(ctx, pgCfg)

	redisOpts := redis.Options{
		Addr: RedisEndpointURL,
	}
	r := redis_app.InitRedis(ctx, redisOpts)
	_, err := r.Ping(ctx).Result()
	if err != nil {
		log.Err(err)
	}

	log.Info().Msg("starting apollo redis")
	beacon_fetcher.InitFetcherService(ctx, BeaconEndpointURL, r)
	log.Info().Interface("redis conn", r.Conn()).Msg("started redis")

	log.Info().Msg("starting apollo s3 manager")
	apollo_buckets.ApolloS3Manager = auth_startup.NewDigitalOceanS3AuthClient(ctx, authKeysCfg)
	log.Info().Msg("connected apollo s3 manager")

	err = e.Start("0.0.0.0:9000")
	if err != nil {
		log.Err(err)
	}
}

func init() {
	viper.AutomaticEnv()

	ApiCmd.Flags().StringVar(&authKeysCfg.AgePubKey, "age-public-key", "age1n97pswc3uqlgt2un9aqn9v4nqu32egmvjulwqp3pv4algyvvuggqaruxjj", "age public key")
	ApiCmd.Flags().StringVar(&authKeysCfg.AgePrivKey, "age-private-key", "", "age private key")
	ApiCmd.Flags().StringVar(&authKeysCfg.SpacesKey, "do-spaces-key", "", "do s3 spaces key")
	ApiCmd.Flags().StringVar(&authKeysCfg.SpacesPrivKey, "do-spaces-private-key", "", "do s3 spaces private key")
	ApiCmd.Flags().StringVar(&RedisEndpointURL, "redis-endpoint", "eth-indexer-redis-headless:6379", "redis endpoint url")
	ApiCmd.Flags().StringVar(&env, "env", "local", "environment")
}

// ApiCmd represents the base command when called without any subcommands
var ApiCmd = &cobra.Command{
	Use:   "eth2-indexer",
	Short: "A Reporting Focused Indexer for Ethereum 2.0/Consensus Layer",
	Run: func(cmd *cobra.Command, args []string) {
		Api()
	},
}
