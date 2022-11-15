package base

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/test"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type FileBaseTestSuite struct {
	test.AutoGenBaseTestSuiteBase
}

func (s *FileBaseTestSuite) TestCodeGen() {
	fb := FileComponentBaseElements{}
	p := structs.Path{PackageName: "base", FnIn: "base_file_example.go"}

	f := fb.GenerateFileShell(p)
	s.Assert().NotEmpty(f)
	err := f.Save(p.FnIn)
	s.Assert().Nil(err)

	s.Cleanup = false
	if s.Cleanup {
		s.DeleteFile(p.FnIn)
	}
}

func TestFileBaseTestSuite(t *testing.T) {
	suite.Run(t, new(FileBaseTestSuite))
}
