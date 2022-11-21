package ecdsa_signer

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ECDSATestSuite struct {
	suite.Suite
}

func (s *ECDSATestSuite) TestNewSignerInit() {
	pkHexString := NewEcdsaPkHexString()
	es, err := CreateEcdsaSignerFromPk(pkHexString)
	s.Assert().Nil(err)
	s.Assert().NotNil(es.Account)
}

func TestECDSATestSuite(t *testing.T) {
	suite.Run(t, new(ECDSATestSuite))
}
