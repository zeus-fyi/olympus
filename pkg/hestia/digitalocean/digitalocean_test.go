package hestia_digitalocean

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_nodes "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/nodes"
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

	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}
func (s *DigitalOceanTestSuite) TestListSizes() {
	sizes, err := s.do.GetSizes(ctx)
	s.Require().NoError(err)
	s.Require().NotEmpty(sizes)
	fmt.Println(sizes)
}

func (s *DigitalOceanTestSuite) TestInsertSizes() {
	sizes, err := s.do.GetSizes(ctx)
	s.Require().NoError(err)
	s.Require().NotEmpty(sizes)
	nodes := hestia_autogen_bases.NodesSlice{}
	for _, size := range sizes {
		for _, reg := range size.Regions {
			dbSize := hestia_autogen_bases.Nodes{}
			dbSize.Slug = size.Slug
			dbSize.Disk = size.Disk
			dbSize.PriceHourly = size.PriceHourly
			dbSize.CloudProvider = "do"
			dbSize.Vcpus = size.Vcpus
			dbSize.Region = reg
			dbSize.PriceMonthly = size.PriceMonthly
			dbSize.Memory = size.Memory
			dbSize.Description = size.Description
			nodes = append(nodes, dbSize)
		}
	}

	err = hestia_nodes.InsertNodes(ctx, nodes)
	s.Require().NoError(err)
	fmt.Println(nodes)
}

func TestDigitalOceanTestSuite(t *testing.T) {
	suite.Run(t, new(DigitalOceanTestSuite))
}
