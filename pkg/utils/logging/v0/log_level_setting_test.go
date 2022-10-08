package v0

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
}

func (s *LoggerTestSuite) TestSetLevel() {
	l := LibV0{}

	level := l.SetLoggerLevel(zerolog.Level(0))
	s.Assert().NotEmpty(level)
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
