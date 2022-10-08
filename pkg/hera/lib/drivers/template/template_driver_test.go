package template

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type ToJenDriverTestSuite struct {
	base.TestSuite
}

func (s *ToJenDriverTestSuite) SetupTest() {
}

func pathCreationForTemplateTest() structs.Path {
	dirIn := "models"
	dirOut := "autogen_template_preview"
	fn := "model_template.go"
	pkgName := "models"
	env := "test"
	pathIn := p.NewFullPathDefinition(env, pkgName, dirIn, dirOut, fn)
	return pathIn
}
func (s *ToJenDriverTestSuite) TestSingleFileTemplateGeneration() {
	pathIn := pathCreationForTemplateTest()
	s.Require().Nil(CreateJenFile(pathIn))
	pathIn.DirOut = "autogen"
	//s.Require().Nil(p.CleanUpPaths(pathIn))
}

func (s *ToJenDriverTestSuite) TestDirectoryTemplateGeneration() {
	pathIn := pathCreationForTemplateTest()
	s.Require().Nil(CreateJenFilesFromDir(pathIn))
	pathIn.DirOut = "autogen"
	//s.Require().Nil(p.CleanUpPaths(pathIn))
}

func TestToJenDriverTestSuite(t *testing.T) {
	suite.Run(t, new(ToJenDriverTestSuite))
}
