package v1_iris

import (
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/labstack/echo/v4"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	"github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

func (s *IrisV1TestSuite) TestAnvilSessionLock() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	iris_redis.InitLocalTestProductionRedisIrisCache(ctx)

	InitV1Routes(s.E)
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = s.E.Start(":9010")
	}()

	groupName := "ethereum-mainnet"
	irisClient := resty_base.GetBaseRestyClient("http://localhost:9010", s.Tc.ProductionLocalTemporalBearerToken)
	irisClient.Header.Set("Content-Type", "application/json")
	irisClient.Header.Set(iris_programmable_proxy.RouteGroupHeader, groupName)
	irisClient.Header.Set(AnvilSessionLockHeader, "Zeus-Test")
	payload := `{
		"jsonrpc": "2.0",
		"method": "eth_getBlockByNumber",
		"params": ["latest", true],
		"id": 1
	}`

	resp, err := irisClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("/v1/router")

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().NotEqual(http.StatusNotFound, resp.StatusCode())
}

func (s *IrisV1TestSuite) TestSwitch() {
	b := echo.Map{
		"jsonrpc": "2.0",
		"method":  "hardhat_reset",
	}
	method, ok := b["method"].(string)
	if ok && strings.HasPrefix(method, "hardhat_") {
		b["method"] = replacePrefix(method, "hardhat_", "anvil_")
	}

	s.Require().Equal("anvil_reset", b["method"])
}
