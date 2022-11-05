package workflows

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type TemporalWorkflowsTestSuite struct {
	base.TestSuite
}

func (s *TemporalWorkflowsTestSuite) TestCreateWorkflow() {
	// TODO
}

func TestTemporalWorkflowsTestSuite(t *testing.T) {
	suite.Run(t, new(TemporalWorkflowsTestSuite))
}
