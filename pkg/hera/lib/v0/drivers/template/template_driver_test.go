package template_driver

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type ToJenDriverTestSuite struct {
	test_suites_base.TestSuite
}

func (s *ToJenDriverTestSuite) SetupTest() {
}

func (s *ToJenDriverTestSuite) TestTemplateCreation() {
}

func TestToJenDriverTestSuite(t *testing.T) {
	suite.Run(t, new(ToJenDriverTestSuite))
}
