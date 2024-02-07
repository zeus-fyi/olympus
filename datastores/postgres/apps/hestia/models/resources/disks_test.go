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

func (s *DisksTestSuite) TestInsertDiskOvh() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	// hourlyDiskCost = 0.01643835616
	// monthly = 12
	disk := hestia_autogen_bases.Disks{
		PriceMonthly:  12,
		PriceHourly:   0.01643835616,
		Region:        "us-west-or-1",
		CloudProvider: "ovh",
		Description:   "OvhCloud Block Storage High Speed SSD",
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
	gbInGi *= 0.1
	fmt.Println(gbInGi, gbInGi/730)
	disk := hestia_autogen_bases.Disks{
		PriceMonthly:  gbInGi,
		PriceHourly:   gbInGi / 730,
		Region:        "us-east-2",
		CloudProvider: "aws",
		Description:   "EBS gp2 Block Storage SSD",
		Type:          "ssd",
		SubType:       "gp2",
		DiskSize:      100,
		DiskUnits:     "Gi",
	}
	err := InsertDisk(ctx, disk)
	s.Require().NoError(err)
}

/*
TODO
General Purpose SSD (gp3) - IOPS	3,000 IOPS free and $0.005/provisioned IOPS-month over 3,000
General Purpose SSD (gp3) - Throughput	125 MB/s free and $0.040/provisioned MB/s-month over 125
*/

func TestDisksTestSuite(t *testing.T) {
	suite.Run(t, new(DisksTestSuite))
}
