package ecdsa_signer

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type ECDSATestSuite struct {
	test_suites.EncryptionTestSuite
}

func (s *ECDSATestSuite) TestLocalEcsdaKey() {
	pkHexString := s.Tc.LocalEcsdaTestPkey
	s.Assert().NotEmpty(pkHexString)
	es, err := CreateEcdsaSignerFromPk(pkHexString)
	s.Assert().Nil(err)
	s.Assert().NotNil(es.Account)
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
