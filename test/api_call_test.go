package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type DriverTestSuite struct {
	test_suites.TemporalTestSuite
}

func (t *DriverTestSuite) TestDriveDeployProdDemoChartApiCall() {
	err := DeployDemoProdChartApiCall()
	t.Require().Nil(err)
}

func (t *DriverTestSuite) TestDriveCreateInternalProdNsApiCall() {
	err := CreateInternalProdNs()
	t.Require().Nil(err)
}

func (t *DriverTestSuite) TestDriveProdDemoCreateTopologyWithChart() {
	err := CreateDemoChartApiCall()
	t.Require().Nil(err)
}

func (t *DriverTestSuite) TestDriveUpdateDeploymentStatus() {
	err := UpdateDeploymentStatusApiCall()
	t.Require().Nil(err)
}

func (t *DriverTestSuite) TestDriveReadTopologiesMetadata() {
	err := ReadTopologiesMetadataAPICall()
	t.Require().Nil(err)
}

func (t *DriverTestSuite) TestDriveReadChart() {
	err := ReadChartAPICall()
	t.Require().Nil(err)
}

func (t *DriverTestSuite) TestDriveCreateTopologyWithChart() {
	err := CreateChartApiCall()
	t.Require().Nil(err)
}

func TestDriverTestSuite(t *testing.T) {
	suite.Run(t, new(DriverTestSuite))
}
