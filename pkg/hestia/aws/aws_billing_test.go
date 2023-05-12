package hestia_eks_aws

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type AwsPricingClientTestSuite struct {
	test_suites_base.TestSuite
	pc AwsPricing
}

func (s *AwsPricingClientTestSuite) SetupTest() {
	s.InitLocalConfigs()
	eksCreds := EksCredentials{
		Region:       "us-east-1",
		AccessKey:    s.Tc.AwsAccessKeyEks,
		AccessSecret: s.Tc.AwsSecretKeyEks,
	}
	p, err := InitPricingClient(ctx, eksCreds)
	s.Require().NoError(err)
	s.pc = p
	s.Require().NotNil(s.pc.Client)
}

func (s *AwsPricingClientTestSuite) TestGetEC2Products() {
	err := s.pc.GetAllProducts(ctx, UsWest1)
	s.Require().NoError(err)
}
func (s *AwsPricingClientTestSuite) TestGetEC2Product() {
	instanceTypes := []string{
		"t2.nano",
		"t2.micro",
		"t2.small",
		"t2.medium",
		"t2.large",
		"t2.2xlarge",
		"t3.micro",
		"t3.small",
		"t3.medium",
		"t3.xlarge",
	}

	n := hestia_autogen_bases.NodesSlice{}

	fmt.Println(n)
	for _, instanceType := range instanceTypes {
		prices, err := s.pc.GetEC2Product(ctx, UsWest1, instanceType)
		s.Require().NoError(err)

		for _, price := range prices {
			usdCost, timeUnit, perr := price.GetPricePerUnitUSD()
			s.Require().NoError(perr)
			fmt.Printf("Cost: %f %s ", usdCost, timeUnit)
			fmt.Printf("Monthly Cost: %f\n", usdCost*730)
			fmt.Printf("Description: %s\n", price.GetDescription())
			mem, memUnits := price.GetMemoryAndUnits()
			fmt.Printf("Memory: %s, %s\n", mem, memUnits)
			dbSize := hestia_autogen_bases.Nodes{}
			dbSize.Slug = instanceType
			dbSize.Disk = 0
			dbSize.DiskUnits = "GB"
			dbSize.PriceHourly = usdCost
			dbSize.CloudProvider = "aws"
			//dbSize.Vcpus = float64(size.Vcpus)
			//dbSize.Region = reg
			//dbSize.PriceMonthly = size.PriceMonthly
			//dbSize.Memory = size.Memory
			//dbSize.MemoryUnits = "MB"
			//dbSize.Description = size.Description
			//n = append(n, dbSize)
		}

	}

	//err := hestia_compute_resources.InsertNodes(ctx, n)
	//s.Require().NoError(err)
}

func TestAwsPricingClientTestSuite(t *testing.T) {
	suite.Run(t, new(AwsPricingClientTestSuite))
}
