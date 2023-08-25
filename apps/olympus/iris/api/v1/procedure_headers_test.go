package v1_iris

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type IrisV1TestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (s *IrisV1TestSuite) TestHeaders() {
	ph := ProcedureHeaders{
		XAggOp:               "max",
		XAggKey:              "",
		XAggKeyValueDataType: "",
		XAggComp:             "",
		XAggCompDataType:     "",
		XAggFilterPayload:    "",
		XAggFilterFanIn:      nil,
	}
	proc := ph.GetGeneratedProcedure()
	s.NotNil(proc.OrderedSteps)
}

func (s *IrisV1TestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	//iris_redis.InitLocalTestProductionRedisIrisCache(ctx)
}

func TestIrisV1TestSuite(t *testing.T) {
	suite.Run(t, new(IrisV1TestSuite))
}
