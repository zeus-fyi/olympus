package aegis_crypto

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EcdsaTestSuite struct {
	suite.Suite
}

// m/44'/60'/0'/0/0
func (s *EcdsaTestSuite) TestEthWalletGeneration() {
	numWorkers := runtime.NumCPU()

	ag, err := GenAddresses(10000, numWorkers)
	s.Require().Nil(err)

	fmt.Println("Mnemonic: ", ag.Mnemonic)
	fmt.Println("Path Index: ", ag.PathIndex)
	fmt.Println("Address: ", ag.Address)
}

func TestEcdsaTestSuite(t *testing.T) {
	suite.Run(t, new(EcdsaTestSuite))
}
