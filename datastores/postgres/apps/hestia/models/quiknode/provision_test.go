package hestia_quicknode_models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

var ctx = context.Background()

type QuickNodeProvisioningTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *QuickNodeProvisioningTestSuite) TestInsertProvisionedService() {

}

func (s *QuickNodeProvisioningTestSuite) TestUpdateProvisionedService() {

}

func (s *QuickNodeProvisioningTestSuite) TestDeprovisionedService() {

}

func (s *QuickNodeProvisioningTestSuite) TestDeactivateService() {

}

func TestQuickNodeProvisioningTestSuite(t *testing.T) {
	suite.Run(t, new(QuickNodeProvisioningTestSuite))
}
