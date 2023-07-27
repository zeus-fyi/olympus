package hestia_quiknode_v1_routes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hesta_base_test "github.com/zeus-fyi/olympus/hestia/api/test"
)

type QuickNodeTestSuite struct {
	hesta_base_test.HestiaBaseTestSuite
}

var ctx = context.Background()

func (t *QuickNodeTestSuite) Test() {
}

func TestQuickNodeTestSuite(t *testing.T) {
	suite.Run(t, new(QuickNodeTestSuite))
}
