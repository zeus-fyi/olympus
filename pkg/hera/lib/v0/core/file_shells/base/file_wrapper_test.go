package base

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/primitives"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/core/template_test"
)

type FileBaseTestSuite struct {
	template_test.TemplateTestSuite
}

func (s *FileBaseTestSuite) TestCodeGen() {

	fb := FileComponentBaseElements{}
	fw := primitives.FileWrapper{PackageName: "base", FileName: "base_file_example.go"}

	f := fb.GenerateFileShell(fw)
	s.Assert().NotEmpty(f)
	err := f.Save(fw.FileName)
	s.Assert().Nil(err)

	s.Cleanup = false
	if s.Cleanup {
		s.DeleteFile(fw.FileName)
	}
}

func TestFileBaseTestSuite(t *testing.T) {
	suite.Run(t, new(FileBaseTestSuite))
}
