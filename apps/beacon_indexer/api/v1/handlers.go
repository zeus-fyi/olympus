package v1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/beacon-indexer/beacon_indexer/beacon_fetcher"
	"github.com/zeus-fyi/olympus/pkg/logging"
	"github.com/zeus-fyi/olympus/pkg/utils/strings"
)

func Health(c echo.Context) error {
	return c.String(http.StatusOK, "Healthy")
}

func SetLogLevel(c echo.Context) error {
	level := c.Param("level")
	return c.String(http.StatusOK, "Set logging level to : "+logging.SetLoggerLevel(level))
}

func SetNewValidatorBatchSize(c echo.Context) error {
	batchSize := c.Param("batchSize")
	querySize := strings.IntStringParser(batchSize)
	beacon_fetcher.NewValidatorBatchSize = querySize
	return c.String(http.StatusOK, "SetNewValidatorBatchSize: "+batchSize)
}

func SetNewValidatorBalanceBatchSize(c echo.Context) error {
	batchSize := c.Param("batchSize")
	querySize := strings.IntStringParser(batchSize)
	beacon_fetcher.NewValidatorBalancesBatchSize = querySize
	return c.String(http.StatusOK, "SetNewValidatorBatchSize: "+batchSize)
}

func SetNewValidatorBalanceFetchTimeout(c echo.Context) error {
	seconds := c.Param("seconds")
	beacon_fetcher.NewValidatorBalancesTimeout = time.Duration(strings.Int64StringParser(seconds)) * time.Second
	return c.String(http.StatusOK, "SetNewValidatorBalanceFetchTimeout: "+seconds)
}
