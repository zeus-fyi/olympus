package beacon_api

import (
	"fmt"
	"testing"

	"bitbucket.org/zeus/eth-indexer/configs"
	"github.com/stretchr/testify/suite"

	"bitbucket.org/zeus/eth-indexer/pkg/test_utils"
)

type PGTestSuite struct {
	suite.Suite
	tc configs.TestContainer
}

func (s *PGTestSuite) SetupTest() {
	s.tc = test_utils.InitLocalConfigs()
}

func (s *PGTestSuite) TestGetValidatorsByState() {
	slot := "4195194"
	r := GetValidatorsByState(s.tc.BEACON_NODE_INFURA, slot)
	s.Assert().Nil(r.Err)
	s.Assert().NotEmpty(r.Body)
	fmt.Println(r.Body)
}

func (s *PGTestSuite) TestGetBlockByID() {
	r := GetBlockByID(s.tc.BEACON_NODE_INFURA, "100000")
	s.Assert().Nil(r.Err)
	s.Assert().NotEmpty(r.Body)
	fmt.Println(r.Body)
}

func TestPGTestSuite(t *testing.T) {
	suite.Run(t, new(PGTestSuite))
}
