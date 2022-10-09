package code_driver

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type DriverTestSuite struct {
	base.TestSuite
}

func (s *DriverTestSuite) SetupTest() {
}

func TestDriverTestSuite(t *testing.T) {
	suite.Run(t, new(DriverTestSuite))
}
