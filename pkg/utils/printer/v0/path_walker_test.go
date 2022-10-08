package v0

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
)

type PathWalkerTestSuite struct {
	suite.Suite
}

func (s *PathWalkerTestSuite) TestPathWalker() {
	l := Lib{}
	p := structs.Path{DirIn: "."}
	paths := l.BuildPathsFromDirInPath(p, ".go")
	s.Assert().NotEmpty(paths)
}

func TestPathWalkerTestSuite(t *testing.T) {
	suite.Run(t, new(PathWalkerTestSuite))
}
