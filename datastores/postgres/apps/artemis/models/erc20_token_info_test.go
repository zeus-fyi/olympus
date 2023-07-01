package artemis_validator_service_groups_models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type ERC20TokenInfoTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *ERC20TokenInfoTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
}

func (s *ERC20TokenInfoTestSuite) TestInsert() {
	tokenInfo := artemis_autogen_bases.Erc20TokenInfo{
		Address:           "0x",
		ProtocolNetworkID: 1,
		BalanceOfSlotNum:  0,
	}
	err := InsertERC20TokenInfo(ctx, tokenInfo)
	s.Require().Nil(err)
}

func (s *ERC20TokenInfoTestSuite) TestSelectAll() {
	tokens, _, err := SelectERC20Tokens(ctx)
	s.Require().Nil(err)
	fmt.Println(len(tokens))
}

func (s *ERC20TokenInfoTestSuite) TestSelect() {
	tokenInfo := artemis_autogen_bases.Erc20TokenInfo{}
	slotNum, err := SelectERC20TokenInfo(ctx, tokenInfo)
	s.Require().Nil(err)
	s.Assert().Equal(-1, slotNum)

	tokenInfo = artemis_autogen_bases.Erc20TokenInfo{
		Address:           "0x",
		ProtocolNetworkID: 1,
	}
	slotNum, err = SelectERC20TokenInfo(ctx, tokenInfo)
	s.Require().Nil(err)
	s.Assert().Equal(0, slotNum)

}

func TestERC20TokenInfoTestSuite(t *testing.T) {
	suite.Run(t, new(ERC20TokenInfoTestSuite))
}
