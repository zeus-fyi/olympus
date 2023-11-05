package v1_promql

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	mev_promql "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/promql"
)

type PromQLRequest struct{}

func PromQLRequestHandler(method string) func(c echo.Context) error {
	return func(c echo.Context) error {
		request := new(PromQLRequest)
		if err := c.Bind(request); err != nil {
			return err
		}
		return request.PromQLRequest(c, method)
	}
}

func (p *PromQLRequest) PromQLRequest(c echo.Context, method string) error {
	switch method {
	case "top_k_tokens":
		timeNow := time.Now().UTC()
		window := v1.Range{
			Start: timeNow.Add(-time.Minute * 60),
			End:   time.Now().UTC(),
			Step:  time.Minute,
		}
		resp, err := mev_promql.ProxyMevPromQL.GetTopTokens(c.Request().Context(), window)
		if err != nil {
			return err
		}
		tokens := make([]string, len(resp))
		for i, val := range resp {
			tokens[i] = val.Metric.In
		}
		return c.JSON(http.StatusOK, tokens)
	default:
	}
	return c.JSON(http.StatusNotImplemented, nil)
}
