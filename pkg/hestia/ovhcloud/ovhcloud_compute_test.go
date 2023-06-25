package hestia_ovhcloud

import "fmt"

func (s *OvhCloudTestSuite) TestInsertSizes() {
	sizes, err := s.o.GetSizes(ctx, OvhServiceName, OvhRegionUsWestOr1ENUM)
	s.Require().NoError(err)
	s.Require().NotEmpty(sizes)

	for _, sa := range sizes {
		fmt.Println(sa)
	}

	//n := hestia_autogen_bases.NodesSlice{}
	////for _, size := range sizes {
	////	for _, reg := range size.Regions {
	////		dbSize := hestia_autogen_bases.Nodes{}
	////		dbSize.Slug = size.Slug
	////		dbSize.Disk = size.Disk
	////		dbSize.DiskUnits = "GB"
	////		dbSize.PriceHourly = size.PriceHourly
	////		dbSize.CloudProvider = "do"
	////		dbSize.Vcpus = float64(size.Vcpus)
	////		dbSize.Region = reg
	////		dbSize.PriceMonthly = size.PriceMonthly
	////		dbSize.Memory = size.Memory
	////		dbSize.MemoryUnits = "MB"
	////		dbSize.Description = size.Description
	////		n = append(n, dbSize)
	////	}
	////}
	//
	//err = hestia_compute_resources.InsertNodes(ctx, n)
	//s.Require().NoError(err)
	//fmt.Println(n)
}
