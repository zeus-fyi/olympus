package hestia_digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type DomainsTestSuite struct {
	DigitalOceanTestSuite
}

var ctx = context.Background()

func (s *DomainsTestSuite) TestCreateSubdomain() {
	cctx := zeus_common_types.CloudCtxNs{Namespace: "b9e8864a-a829-4d1c-9c36-56e21bb95eab"}
	dr, err := s.do.CreateDomain(ctx, cctx)
	s.Require().NoError(err)
	s.Require().NotNil(dr)
}

func (s *DomainsTestSuite) TestRemoveSubDomainARecord() {
	cctx := zeus_common_types.CloudCtxNs{Namespace: "b9e8864a-a829-4d1c-9c36-56e21bb95eab"}
	err := s.do.RemoveSubDomainARecord(ctx, cctx)
	s.Require().NoError(err)
}

func (s *DomainsTestSuite) TestListCnames() {
	dl, _, err := s.do.Domains.RecordsByType(ctx, "zeus.fyi", "CNAME", &godo.ListOptions{})
	s.Require().NoError(err)
	for _, dn := range dl {
		fmt.Println(dn.Name, dn.Data)
	}
}
func TestDomainsTestSuite(t *testing.T) {
	suite.Run(t, new(DomainsTestSuite))
}
