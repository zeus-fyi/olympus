package apollo_ethereum_alerts

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/pretty"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type ApolloEthereumAlertsTestSuite struct {
	pc ApolloEthereumAlerts
	test_suites_base.TestSuite
}

func (t *ApolloEthereumAlertsTestSuite) SetupTest() {
	t.InitLocalConfigs()
	t.pc = InitLocalApolloEthereumAlerts(ctx, t.Tc.PagerDutyApiKey, t.Tc.PagerDutyRoutingKey)
}

func (t *ApolloEthereumAlertsTestSuite) TestSlashingQueryRangePromQL() {
	_, err := t.pc.SlashingAlertTrigger(ctx)
	t.Require().NoError(err)
}

func (t *ApolloEthereumAlertsTestSuite) TestQueryRangePromQL() {
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

	mv := r.(model.Matrix)
	fmt.Println(mv)

	b, err := json.Marshal(r)
	t.Require().NoError(err)
	requestJSON := pretty.Pretty(b)
	requestJSON = pretty.Color(requestJSON, pretty.TerminalStyle)
	fmt.Println(string(requestJSON))
}

func TestApolloEthereumAlertsTestSuite(t *testing.T) {
	suite.Run(t, new(ApolloEthereumAlertsTestSuite))
}
