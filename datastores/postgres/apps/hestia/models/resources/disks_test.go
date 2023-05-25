package hestia_compute_resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type DisksTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *DisksTestSuite) TestInsertDisk() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	disk := hestia_autogen_bases.Disks{
		PriceMonthly:  10,
		PriceHourly:   0.015,
		Region:        "nyc1",
		CloudProvider: "do",
		Description:   "Digital Ocean Block Storage SSD",
		Type:          "ssd",
		DiskSize:      100,
		DiskUnits:     "Gi",
	}
	err := InsertDisk(ctx, disk)
	s.Require().NoError(err)
}

func (s *DisksTestSuite) TestInsertDiskGcp() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	disk := hestia_autogen_bases.Disks{
		PriceMonthly:  17,
		PriceHourly:   0.02329,
		Region:        "us-central1",
		CloudProvider: "gcp",
		Description:   "GCP Block Storage SSD",
		Type:          "ssd",
		DiskSize:      100,
		DiskUnits:     "Gi",
	}
	err := InsertDisk(ctx, disk)
	s.Require().NoError(err)
}

func (s *DisksTestSuite) TestInsertDiskAws() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	gbInGi := 107.374
	gbInGi *= 0.12
	fmt.Println(gbInGi, gbInGi/730)
	disk := hestia_autogen_bases.Disks{
		PriceMonthly:  gbInGi,
		PriceHourly:   gbInGi / 730,
		Region:        "us-west-1",
		CloudProvider: "aws",
		Description:   "EBS gp2 Block Storage SSD",
		Type:          "ssd",
		DiskSize:      100,
		DiskUnits:     "Gi",
	}
	err := InsertDisk(ctx, disk)
	s.Require().NoError(err)
}
func TestDisksTestSuite(t *testing.T) {
	suite.Run(t, new(DisksTestSuite))
}
