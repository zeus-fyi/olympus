package conversions_test

import (
	"os"
	"path"
	"runtime"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
)

var PgTestDB apps.Db

type ConversionsTestSuite struct {
	test_suites.PGTestSuite
	Yr            transformations.YamlReader
	TestDirectory string
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

func (s *ConversionsTestSuite) SetupTest() {
	s.TestDirectory = ForceDirToCallerLocation()
	s.Yr = transformations.YamlReader{}
	s.InitLocalConfigs()
	s.SetupPGConn()
}
