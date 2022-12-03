package code_driver

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type DriverTestSuite struct {
	test_suites_base.TestSuite
}

func (s *DriverTestSuite) SetupTest() {
}

func TestDriverTestSuite(t *testing.T) {
	suite.Run(t, new(DriverTestSuite))
}
