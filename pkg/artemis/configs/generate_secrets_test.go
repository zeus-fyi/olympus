package artemis_network_cfgs

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ArtemisGenerateSecretsTestSuite struct {
	suite.Suite
}

func (s *ArtemisGenerateSecretsTestSuite) TestSecretsGen() {

}

func TestArtemisGenerateSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisGenerateSecretsTestSuite))
}
