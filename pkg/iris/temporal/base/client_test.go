package temporal_base

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type TemporalClientTestSuite struct {
	base.TestSuite
}

func (s *TemporalClientTestSuite) SetupTest() {
}

func (s *TemporalClientTestSuite) TestCreateClient() {
	s.InitLocalConfigs()
	tc, err := NewTemporalClient(s.Tc.DevTemporalAuth)
	s.Require().Nil(err)
	err = tc.ConnectTemporalClient()
	s.Require().Nil(err)
	defer tc.Close()
}

func TestTemporalClientTestSuite(t *testing.T) {
	suite.Run(t, new(TemporalClientTestSuite))
}
