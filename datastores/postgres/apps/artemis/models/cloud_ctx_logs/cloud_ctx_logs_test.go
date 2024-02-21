package cloud_ctx_logs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type CloudCtxLogsTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CloudCtxLogsTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.Pg.InitPG(context.Background(), s.Tc.LocalDbPgconn)
}

var (
	ctx = context.Background()
)

func (s *CloudCtxLogsTestSuite) TestInsertCloudCtxNsLog() {
	cl := CloudCtxNsLogs{
		LogID:           0,
		OrchestrationID: 1679548290001220864,
		Status:          "info",
		Msg:             "test log message",
		Ou:              s.Ou,
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			CloudProvider: "do",
			Region:        "sfo3",
			Context:       "do-sfo3-dev-do-sfo3-zeus",
			Namespace:     "ethereum",
		},
	}
	err := InsertCloudCtxNsLog(ctx, &cl)
	s.Require().Nil(err)
	s.Require().NotZero(cl.LogID)
}

func (s *CloudCtxLogsTestSuite) TestSelectCloudCtxNsLogs() {
	res, err := SelectCloudCtxNsLogs(ctx, CloudCtxNsLogs{
		CloudCtxNsID: 1668716905827494000,
		Ou:           s.Ou,
	})
	s.Require().Nil(err)
	s.Require().Greaterf(len(res), 1, "expected at least 1 log entry")
}

func TestCloudCtxLogsTestSuite(t *testing.T) {
	suite.Run(t, new(CloudCtxLogsTestSuite))
}
