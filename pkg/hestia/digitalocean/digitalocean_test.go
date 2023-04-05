package hestia_digitalocean

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type DigitalOceanTestSuite struct {
	test_suites_base.TestSuite
	do DigitalOcean
}

func (s *DigitalOceanTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.do = InitDoClient(ctx, s.Tc.DigitalOceanAPIKey)
	s.Require().NotNil(s.do.Client)
}
func (s *DigitalOceanTestSuite) TestListSizes() {
	sizes, err := s.do.GetSizes(ctx)
	s.Require().NoError(err)
	s.Require().NotEmpty(sizes)
	fmt.Println(sizes)
}

func TestDigitalOceanTestSuite(t *testing.T) {
	suite.Run(t, new(DigitalOceanTestSuite))
}
