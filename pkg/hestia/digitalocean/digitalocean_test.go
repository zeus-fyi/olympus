package hestia_digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	"k8s.io/apimachinery/pkg/api/resource"
)

type DigitalOceanTestSuite struct {
	test_suites_base.TestSuite
	do DigitalOcean
}

func (s *DigitalOceanTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.do = InitDoClient(ctx, s.Tc.DigitalOceanAPIKey)
	s.Require().NotNil(s.do.Client)
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	//apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
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
	n := hestia_autogen_bases.NodesSlice{}
	for _, size := range sizes {
		for _, reg := range size.Regions {
			dbSize := hestia_autogen_bases.Nodes{}
			dbSize.Slug = size.Slug
			dbSize.Disk = size.Disk
			dbSize.DiskUnits = "GB"
			dbSize.PriceHourly = size.PriceHourly
			dbSize.CloudProvider = "do"
			dbSize.Vcpus = float64(size.Vcpus)
			dbSize.Region = reg
			dbSize.PriceMonthly = size.PriceMonthly
			dbSize.Memory = size.Memory
			dbSize.MemoryUnits = "MB"
			dbSize.Description = size.Description
			n = append(n, dbSize)
		}
	}

	err = hestia_compute_resources.InsertNodes(ctx, n)
	s.Require().NoError(err)
	fmt.Println(n)
}

func (s *DigitalOceanTestSuite) TestQuantity() {

	qty, err := digitalOceanBlockStorageBillingUnits(ctx, "10Gi")
	s.Require().NoError(err)
	s.Require().Equal(0.1, qty)
}

func digitalOceanBlockStorageBillingUnits(ctx context.Context, qtyString string) (float64, error) {
	r, err := resource.ParseQuantity(qtyString)
	if err != nil {
		return 0, err
	}
	rawValue := r.Value()
	q := resource.NewQuantity(rawValue, resource.BinarySI)
	q.ScaledValue(resource.Giga)
	miValue := float64(q.AsDec().UnscaledBig().Int64() / (1024 * 1024 * 1024))
	billableUnits := miValue / 100
	return billableUnits, nil
}

func TestDigitalOceanTestSuite(t *testing.T) {
	suite.Run(t, new(DigitalOceanTestSuite))
}
