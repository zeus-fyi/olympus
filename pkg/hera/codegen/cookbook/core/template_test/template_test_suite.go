package template_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/printer"
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
)

type TemplateTestSuite struct {
	suite.Suite
	Cleanup bool
}

func (s *TemplateTestSuite) SetupTest() {
	s.Cleanup = true
}

func (s *TemplateTestSuite) DeleteFile(fn string) {
	path := structs.Path{Fn: fn}
	p := printer.Printer{}
	s.Require().Nil(p.DeleteFile(path))
}

func TestTemplateTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateTestSuite))
}
