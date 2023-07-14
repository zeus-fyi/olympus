package iris_models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

var ctx = context.Background()

type IrisTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func TestIrisTestSuite(t *testing.T) {
	suite.Run(t, new(IrisTestSuite))
}
