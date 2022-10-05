package template_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/printer"
)

type TemplateTestSuite struct {
	suite.Suite
	Cleanup bool
}

func (s *TemplateTestSuite) SetupTest() {
	s.Cleanup = true
}

func (s *TemplateTestSuite) DeleteFile(fn string) {
	err := printer.DeleteFile(fn)
	s.Require().Nil(err)
}

func TestTemplateTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateTestSuite))
}
