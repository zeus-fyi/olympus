package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type AutoGenBaseTestSuiteBase struct {
	test_suites_base.TestSuite
	Cleanup bool
}

func (s *AutoGenBaseTestSuiteBase) SetupTest() {
	s.Cleanup = true
}

func (s *AutoGenBaseTestSuiteBase) DeleteFile(fn string) {
	path := filepaths.Path{FnIn: fn}
	p := file_io.FileIO{}
	s.Require().Nil(p.DeleteFile(path))
}

func TestAutoGenBaseTestSuiteBase(t *testing.T) {
	suite.Run(t, new(AutoGenBaseTestSuiteBase))
}
