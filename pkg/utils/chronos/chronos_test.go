package chronos

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ChronosTestSuite struct {
	suite.Suite
}

func (s *ChronosTestSuite) TestLib0() {
	c := Chronos{}
	s.Require().NotEmpty(c.v0.UnixTimeStampNow())
}

func TestChronosTestSuite(t *testing.T) {
	suite.Run(t, new(ChronosTestSuite))
}
