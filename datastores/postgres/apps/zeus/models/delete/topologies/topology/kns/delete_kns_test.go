package delete_kns

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type CreateKnsTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *CreateKnsTestSuite) TestDeleteKns() {

	newKns := kns.NewKns()
	newKns.CloudProvider = "do"
	newKns.Region = "sfo3"
	newKns.Context = "context"
	newKns.Env = "test"
	newKns.Namespace = "testnamespace"
	newKns.TopologyID = 1668387026225673081

	ctx := context.Background()
	err := DeleteKns(ctx, &newKns)
	s.Require().Nil(err)
}

func TestCreateKnsTestSuite(t *testing.T) {
	suite.Run(t, new(CreateKnsTestSuite))
}
