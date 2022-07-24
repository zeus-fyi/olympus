package logging

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
}

func (s *LoggerTestSuite) TestSetLevel() {
	level := SetLoggerLevel("0")
	s.Assert().NotEmpty(level)
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
