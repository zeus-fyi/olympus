package paths

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

type PathWalkerTestSuite struct {
	suite.Suite
}

func (s *PathWalkerTestSuite) TestPathWalker() {
	l := PathLib{}
	p := filepaths.Path{DirIn: "."}
	paths := l.BuildPathsFromDirInPath(p, ".go")
	s.Assert().NotEmpty(paths)
}

func TestPathWalkerTestSuite(t *testing.T) {
	suite.Run(t, new(PathWalkerTestSuite))
}
