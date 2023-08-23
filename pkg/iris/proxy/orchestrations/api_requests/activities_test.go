package iris_api_requests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type IrisActivitiesTestSuite struct {
	test_suites_base.TestSuite
}

func (s *IrisActivitiesTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
}

func (s *IrisActivitiesTestSuite) TestBroadcastETL() {
	bc := NewIrisApiRequestsActivities()

	timeOut := time.Second * 10
	pr := &ApiProxyRequest{}
	routes := []string{"1", "2", "3"}
	resp, err := bc.BroadcastETLRequest(ctx, pr, routes, timeOut)
	s.NoError(err)
	s.NotNil(resp)
}

func TestIrisActivitiesTestSuite(t *testing.T) {
	suite.Run(t, new(IrisActivitiesTestSuite))
}
