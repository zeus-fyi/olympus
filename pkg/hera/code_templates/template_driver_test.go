package code_templates

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type ToJenDriverTestSuite struct {
	base.TestSuite
	cleanUp bool
}

func (s *ToJenDriverTestSuite) SetupTest() {
	ForceDirToCallerLocation()
	s.cleanUp = false
}

func pathCreationForTemplateTest() structs.Path {
	dirIn := "models"
	fn := "model_template.go"
	pkgName := "models"
	pathIn := p.NewPkgPath(pkgName, dirIn, fn)
	return pathIn
}
func (s *ToJenDriverTestSuite) TestTemplateGeneration() {
	pathIn := pathCreationForTemplateTest()
	pathOut := pathIn
	pathOut.Dir = "autogen"
	s.Require().Nil(CreateJenFile(pathIn, pathOut))

	if s.cleanUp {
		s.Require().Nil(p.CleanUpPaths(pathOut))
	}
}

func TestToJenDriverTestSuite(t *testing.T) {
	suite.Run(t, new(ToJenDriverTestSuite))
}

func ForceDirToCallerLocation() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}
