package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type AutoGenBaseTestSuiteBase struct {
	base.TestSuite
	Cleanup bool
}

func (s *AutoGenBaseTestSuiteBase) SetupTest() {
	s.Cleanup = true
}

func (s *AutoGenBaseTestSuiteBase) DeleteFile(fn string) {
	path := structs.Path{FnIn: fn}
	p := file_io.FileIO{}
	s.Require().Nil(p.DeleteFile(path))
}

func TestAutoGenBaseTestSuiteBase(t *testing.T) {
	suite.Run(t, new(AutoGenBaseTestSuiteBase))
}
