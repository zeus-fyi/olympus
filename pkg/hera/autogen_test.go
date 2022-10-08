package hera

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type AutoCodeGenTestSuite struct {
	base.TestSuite
}

func (s *AutoCodeGenTestSuite) SetupTest() {
}

func pathCreationForTemplateTest() structs.Path {
	dirIn := "models"
	dirOut := "template_preview"
	fn := "model_template.go"
	pkgName := "models"
	env := "test"
	pathIn := p.NewFullPathDefinition(env, pkgName, dirIn, dirOut, fn)
	return pathIn
}

func (s *AutoCodeGenTestSuite) TestAutoGen() {
	path := pathCreationForTemplateTest()
	err := CreateTemplate(path)
	s.Require().Nil(err)
}

func (s *AutoCodeGenTestSuite) TestAutoGenDir() {
	path := pathCreationForTemplateTest()
	CreateTemplatesInPath(path)
}

func TestAutoCodeGenTestSuite(t *testing.T) {
	suite.Run(t, new(AutoCodeGenTestSuite))
}
