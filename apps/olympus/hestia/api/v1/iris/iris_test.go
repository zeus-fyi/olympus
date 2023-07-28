package hestia_iris_v1_routes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hesta_base_test "github.com/zeus-fyi/olympus/hestia/api/test"
)

type IrisTestSuite struct {
	hesta_base_test.HestiaBaseTestSuite
}

var ctx = context.Background()

func (t *IrisTestSuite) TestCreateOrgRoute() {

}

func TestIrisTestSuite(t *testing.T) {
	suite.Run(t, new(IrisTestSuite))
}
