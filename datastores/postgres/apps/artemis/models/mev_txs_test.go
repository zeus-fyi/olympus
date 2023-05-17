package artemis_validator_service_groups_models

import (
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type MevTxTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func TestMevTxTestSuite(t *testing.T) {
	suite.Run(t, new(MevTxTestSuite))
}
