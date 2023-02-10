package aegis_hydra

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type AegisHydraTestSuite struct {
	test_suites_base.CoreTestSuite
}

func (s *AegisHydraTestSuite) TestHydraAssignmentsFetch() {
	// TODO
}

func TestAegisHydraTestSuite(t *testing.T) {
	suite.Run(t, new(AegisHydraTestSuite))
}
