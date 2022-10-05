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
	"github.com/zeus-fyi/olympus/datastores/postgres"
	"github.com/zeus-fyi/olympus/datastores/redis_app"
)

var (
	RedisEndpointURL  string
	BeaconEndpointURL string
	PGConnStr         string
)

func Api() {
	// Echo instance
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	e := echo.New()
	e = v1.Routes(e)
	ctx := context.Background()
	postgres.Pg = postgres.Db{}
	MaxConn := int32(10)
	MinConn := int32(3)
	MaxConnLifetime := 15 * time.Minute

	pgCfg := postgres.ConfigChangePG{
		MaxConns:          &MaxConn,
		MinConn:           &MinConn,
		MaxConnLifetime:   &MaxConnLifetime,
		HealthCheckPeriod: nil,
	}
	postgres.Pg.InitPG(ctx, PGConnStr)
	_ = postgres.UpdateConfigPG(ctx, pgCfg)

	redisOpts := redis.Options{
		Addr: RedisEndpointURL,
	}
	r := redis_app.InitRedis(ctx, redisOpts)
	_, err := r.Ping(ctx).Result()
	if err != nil {
		log.Err(err)
	}
	beacon_fetcher.InitFetcherService(ctx, BeaconEndpointURL, r)

	log.Info().Interface("redis conn", r.Conn()).Msg("started redis")
	// Start server

	err = e.Start(":9000")
	if err != nil {
		log.Err(err)
	}
}

func init() {
	viper.AutomaticEnv()
	ApiCmd.Flags().StringVar(&PGConnStr, "postgres-conn-str", "", "postgres connection string")
	ApiCmd.Flags().StringVar(&BeaconEndpointURL, "beacon-endpoint", "", "beacon endpoint url")
	ApiCmd.Flags().StringVar(&RedisEndpointURL, "redis-endpoint", "eth-indexer-redis-headless:6379", "redis endpoint url")
}

// ApiCmd represents the base command when called without any subcommands
var ApiCmd = &cobra.Command{
	Use:   "eth2-indexer",
	Short: "A Reporting Focused Indexer for Ethereum 2.0/Consensus Layer",
	Run: func(cmd *cobra.Command, args []string) {
		Api()
	},
}
