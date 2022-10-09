package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type AutoGenBaseTestSuiteBase struct {
	base.TestSuite
}

func (s *AutoGenBaseTestSuiteBase) SetupTest() {
}

func TestAutoGenBaseTestSuiteBase(t *testing.T) {
	suite.Run(t, new(AutoGenBaseTestSuiteBase))
}
