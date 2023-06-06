package hephaestus_build_actions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

var ctx = context.Background()

type BuildActonsTestSuite struct {
	test_suites_base.TestSuite
}

func (s *BuildActonsTestSuite) SetupTest() {
	s.InitLocalConfigs()
}

func (s *BuildActonsTestSuite) TestBuildActions() {
	dataDir = filepaths.Path{
		PackageName: "",
		DirIn:       "/Users/alex/go/Olympus/olympus/",
		DirOut:      "",
		FnIn:        "",
		FnOut:       "",
		Env:         "",
		FilterFiles: nil,
	}
	appName = "zeus"
	Rebuild()
	Upload(ctx)
}

func TestBuildActonsTestSuite(t *testing.T) {
	suite.Run(t, new(BuildActonsTestSuite))
}
