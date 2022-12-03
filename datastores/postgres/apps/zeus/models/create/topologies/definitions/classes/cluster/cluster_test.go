package create_clusters

import (
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type CreateClustersTestSuite struct {
	h hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (s *CreateClustersTestSuite) TestInsertTopologyState() {

}

func TestCreateClustersTestSuite(t *testing.T) {
	suite.Run(t, new(CreateClustersTestSuite))
}
