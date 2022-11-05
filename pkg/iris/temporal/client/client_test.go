package temporal_client

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type TemporalClientTestSuite struct {
	base.TestSuite
}

func (s *TemporalClientTestSuite) TestCreateClient() {

	// TODO
}

func TestTemporalClientTestSuite(t *testing.T) {
	suite.Run(t, new(TemporalClientTestSuite))
}
