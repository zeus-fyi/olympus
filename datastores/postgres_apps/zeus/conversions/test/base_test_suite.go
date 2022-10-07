package conversions_test

import (
	"github.com/zeus-fyi/olympus/datastores/postgres_apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
)

var PgTestDB postgres_apps.Db

type ConversionsTestSuite struct {
	base.TestSuite
	Yr transformations.YamlReader
}

func (s *ConversionsTestSuite) SetupTest() {
	s.Yr = transformations.YamlReader{}
}
