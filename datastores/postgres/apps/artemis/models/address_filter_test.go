package artemis_validator_service_groups_models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type AddressFilterTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *AddressFilterTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
}

func (s *AddressFilterTestSuite) TestSelect() {
	ab, err := SelectSourceAddresses(ctx, 1)
	s.Require().Nil(err)
	s.NotEmpty(ab)

	for _, addr := range ab {
		fmt.Println(addr)
	}
}

func TestAddressFilterTestSuite(t *testing.T) {
	suite.Run(t, new(AddressFilterTestSuite))
}
