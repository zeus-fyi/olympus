package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type DriverTestSuite struct {
	test_suites.TemporalTestSuite
}

func (t *DriverTestSuite) TestDriveReadChart() {
	err := ReadChartAPICall()
	t.Require().Nil(err)
}

func (t *DriverTestSuite) TestDriveCreateTopologyWithChart() {
	err := CallAPI()
	t.Require().Nil(err)
}

func TestDriverTestSuite(t *testing.T) {
	suite.Run(t, new(DriverTestSuite))
}
