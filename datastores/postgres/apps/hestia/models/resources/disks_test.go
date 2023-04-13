package hestia_compute_resources

import (
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

func TestDisksTestSuite(t *testing.T) {
	suite.Run(t, new(DisksTestSuite))
}
