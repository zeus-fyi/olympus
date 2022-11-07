package deploy_workflow

import (
	"testing"

	"github.com/stretchr/testify/suite"
	temporal_client "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type TemporalWorkflowsTestSuite struct {
	test_suites.TemporalTestSuite
	Temporal temporal_client.TemporalClient
}

func (s *TemporalWorkflowsTestSuite) TestCreateWorkflow() {
}

func TestTemporalWorkflowsTestSuite(t *testing.T) {
	suite.Run(t, new(TemporalWorkflowsTestSuite))
}
