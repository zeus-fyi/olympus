package v1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/admin"
)

type AdminConfigRequest struct {
	AdminConfig
}

type AdminDBConfigRequest struct {
	admin.ConfigChangePG
}

type AdminConfig struct {
	LogLevel *zerolog.Level

	ValidatorBatchSize         *int
	ValidatorBalancesBatchSize *int
	ValidatorBalancesTimeout   *time.Duration
}

type AdminConfigReader struct {
	LogLevel zerolog.Level

	ValidatorBatchSize         int
	ValidatorBalancesBatchSize int
	ValidatorBalancesTimeout   time.Duration
}

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}
