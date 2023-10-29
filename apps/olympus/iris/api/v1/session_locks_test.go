package v1_iris

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

func (s *IrisV1TestSuite) TestAnvilSessionLock() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	p := ProxyRequest{}
	s.E.POST("/v1/anvil", p.ProcessLockedSessionRoute)
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = s.E.Start(":9010")
	}()

	groupName := "ethereum-mainnet"
	irisClient := resty_base.GetBaseRestyClient("http://localhost:9010", s.Tc.ProductionLocalTemporalBearerToken)
	irisClient.Header.Set("Content-Type", "application/json")
	irisClient.Header.Set(iris_programmable_proxy.RouteGroupHeader, groupName)
	irisClient.Header.Set(iris_programmable_proxy_v1_beta.LoadBalancingStrategy, iris_programmable_proxy_v1_beta.Adaptive)
	payload := `{
		"jsonrpc": "2.0",
		"procedure": "eth_maxBlockAggReduce",
		"method": "eth_getBlockByNumber",
		"params": ["latest", true],
		"id": 1
	}`
	resp, err := irisClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("/v1/anvil")

	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().NotEqual(http.StatusNotFound, resp.StatusCode())
}
