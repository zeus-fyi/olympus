package cookbook

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type CookbookTestSuiteBase struct {
	base.TestSuite
}

func (s *CookbookTestSuiteBase) SetupTest() {
	UseCookbookDirectory()
}

func TestCookbookTestSuiteBase(t *testing.T) {
	suite.Run(t, new(CookbookTestSuiteBase))
}
