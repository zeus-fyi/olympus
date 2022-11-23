package containers

import (
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type PodSpecInsertTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (p *PodSpecInsertTestSuite) TestDummyInsertChartSubcomponentSpecPodTemplateContainers() {
}

func TestPodSpecInsertTestSuite(t *testing.T) {
	suite.Run(t, new(PodSpecInsertTestSuite))
}
