package apollo_prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/pretty"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type PrometheusTestSuite struct {
	pc Prometheus
	test_suites_base.TestSuite
}

func (t *PrometheusTestSuite) SetupTest() {
	t.InitLocalConfigs()
	t.pc = NewPrometheusLocalClient(ctx)
}

func (t *PrometheusTestSuite) TestQueryRangePromQL() {
	t.Require().NotEmpty(t.pc)

	timeNow := time.Now().UTC()

	window := v1.Range{
		Start: timeNow.Add(-time.Minute * 60),
		End:   time.Now().UTC(),
		Step:  time.Minute,
	}
	namespace := "athena-beacon-goerli"
	query := fmt.Sprintf("ethereum_beacon_exec_sync_status_block_height{namespace=\"%s\"}", namespace)

	opts := v1.WithTimeout(time.Second * 10)
	r, w, err := t.pc.QueryRange(ctx, query, window, opts)
	fmt.Println(w)
	t.Require().NoError(err)
	t.Assert().NotEmpty(r)
	fmt.Println(r)

	b, err := json.Marshal(r)
	t.Require().NoError(err)
	requestJSON := pretty.Pretty(b)
	requestJSON = pretty.Color(requestJSON, pretty.TerminalStyle)
	fmt.Println(string(requestJSON))
}
func (t *PrometheusTestSuite) TestPromGetRules() {
	t.Require().NotEmpty(t.pc)

	r, err := t.pc.Rules(ctx)
	t.Require().NoError(err)
	t.Assert().NotEmpty(r)
	fmt.Println(r)

	b, err := json.Marshal(r)
	t.Require().NoError(err)
	requestJSON := pretty.Pretty(b)
	requestJSON = pretty.Color(requestJSON, pretty.TerminalStyle)
	fmt.Println(string(requestJSON))
}

func TestPrometheusTestSuite(t *testing.T) {
	suite.Run(t, new(PrometheusTestSuite))
}
