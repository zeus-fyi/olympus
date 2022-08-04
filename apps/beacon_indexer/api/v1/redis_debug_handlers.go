package v1

import (
	"context"
	"net/http"
	"os"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_fetcher"
)

type AdminRedisConfigRequest struct {
	DebugRedis `json:"debug"`
}

func DebugRedisRequestHandler(c echo.Context) error {
	log.Info().Msg("DebugRequestHandler")
	ctx := context.Background()

	request := new(AdminRedisConfigRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	val := os.Getenv(request.OsEnv)
	log.Info().Msgf("logging env var: %s, value: %s", request.OsEnv, val)
	log.Info().Msgf("logging addr: %s", request.Addr)

	if request.UseEnv {
		request.Addr = ""
	}

	opts := redis.Options{
		Addr: request.Addr,
	}

	log.Info().Interface("opts setting: ", opts)
	log.Info().Msgf("logging addr: %s", request.Addr)

	resp, err := beacon_fetcher.Fetcher.Cache.Ping(ctx).Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	log.Info().Interface("DebugRedisRequestHandler ping resp: ", resp)
	return c.JSON(http.StatusOK, resp)
}

func DebugReadRedisRequestHandler(c echo.Context) error {
	log.Info().Msg("DebugRequestHandler")
	ctx := context.Background()

	request := new(AdminRedisConfigRequest)
	if err := c.Bind(request); err != nil {
		return err
	}

	log.Info().Interface("redis setting: ", beacon_fetcher.Fetcher.Cache.Info(ctx))
	resp, err := beacon_fetcher.Fetcher.Cache.Ping(ctx).Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, resp)
}

type DebugRedis struct {
	Addr   string `json:"addr"`
	OsEnv  string `json:"envs,omitempty"`
	UseEnv bool   `json:"enabledEnv,omitempty"`
}
