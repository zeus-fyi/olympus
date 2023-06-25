package hestia_ovhcloud

import (
	"fmt"
	"strings"

	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
)

func (s *OvhCloudTestSuite) TestInsertSizes() {
	sizes, err := s.o.GetSizes(ctx, OvhServiceName, OvhRegionUsWestOr1ENUM)
	s.Require().NoError(err)
	s.Require().NotEmpty(sizes)
	details, err := s.o.GetFlavorDetails(ctx, OvhServiceName, OvhRegionUsWestOr1ENUM)
	s.Require().NoError(err)
	s.Require().NotEmpty(details)
	//for _, sz := range sizes {
	//	fmt.Println(sz.Category)
	//}
	m := make(map[string]Details)
	for _, d := range details {
		m[d.Name] = d
	}

	n := hestia_autogen_bases.NodesSlice{}
	for _, size := range sizes {
		if size.Category == "d" {
			continue
		}
		if strings.HasSuffix(size.Name, "-flex") {
			continue
		}
		flavorDetails := m[size.Name]
		dbSize := hestia_autogen_bases.Nodes{}
		dbSize.Slug = size.Name
		dbSize.Disk = flavorDetails.Disk
		dbSize.DiskUnits = "GB"
		dbSize.PriceHourly = priceMapHourly[size.Name]
		dbSize.CloudProvider = "ovh"
		dbSize.Vcpus = float64(flavorDetails.Vcpus)
		dbSize.Region = OvhRegionUsWestOr1
		dbSize.PriceMonthly = priceMapMonthly[size.Name]
		dbSize.Memory = flavorDetails.Ram
		dbSize.MemoryUnits = "MB"
		dbSize.Description = fmt.Sprintf("%s %d vCPUs %d GB RAM", size.Name, size.VCPUs, size.Ram)
		n = append(n, dbSize)
	}
	err = hestia_compute_resources.InsertNodes(ctx, n)
	s.Require().NoError(err)
	fmt.Println(n)
}
