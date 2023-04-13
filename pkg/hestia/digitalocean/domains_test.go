package hestia_digitalocean

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DomainsTestSuite struct {
	DigitalOceanTestSuite
}

var ctx = context.Background()

func (s *DomainsTestSuite) TestCreateSubdomain() {
	dr, err := s.do.CreateDomain(ctx, "b9e8864a-a829-4d1c-9c36-56e21bb95eab")
	s.Require().NoError(err)
	s.Require().NotNil(dr)
}

func (s *DomainsTestSuite) TestRemoveSubDomainARecord() {
	err := s.do.RemoveSubDomainARecord(ctx, "b9e8864a-a829-4d1c-9c36-56e21bb95eab")
	s.Require().NoError(err)
}

func TestDomainsTestSuite(t *testing.T) {
	suite.Run(t, new(DomainsTestSuite))
}
