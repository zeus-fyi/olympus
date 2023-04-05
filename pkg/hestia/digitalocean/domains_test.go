package hestia_digitalocean

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type DomainsTestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (s *DomainsTestSuite) TestCreateSubdomain() {
	s.InitLocalConfigs()

	InitDoClient(ctx, "token")
}

func TestDomainsTestSuite(t *testing.T) {
	suite.Run(t, new(DomainsTestSuite))
}
