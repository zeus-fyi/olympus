package hestia_eks_aws

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_compute_resources "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/resources"
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
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	instanceTypes := []string{
		//"t2.nano",
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
	for _, instanceType := range instanceTypes {
		prices, err := s.pc.GetEC2Product(ctx, UsWest1, instanceType)
		s.Require().NoError(err)
		fmt.Printf("%s\n", instanceType)
		for _, price := range prices {
			desc := price.GetDescription()
			if !strings.Contains(desc, fmt.Sprintf("per On Demand Linux %s Instance Hour", instanceType)) {
				continue
			}

			trimmed := strings.SplitAfter(desc, "per ")
			fmt.Println(trimmed[1])
			trimmed = strings.SplitAfter(trimmed[1], "Instance")
			desc = strings.TrimSpace(trimmed[0])
			usdCost, timeUnit, perr := price.GetPricePerUnitUSD()
			s.Require().NoError(perr)
			fmt.Printf("Cost: %f %s ", usdCost, timeUnit)
			fmt.Printf("Monthly Cost: %f\n", usdCost*730)
			fmt.Printf("Description: %s\n", desc)
			mem, memUnits := price.GetMemoryAndUnits()
			fmt.Printf("Memory: %s, %s\n", mem, memUnits)
			memInt, merr := strconv.Atoi(mem)
			if merr != nil {
				continue
			}
			s.Require().NoError(merr)

			vcpus, verr := strconv.ParseFloat(price.GetVCpus(), 64)
			s.Require().NoError(verr)
			dbSize := hestia_autogen_bases.Nodes{}
			dbSize.Slug = instanceType
			dbSize.Disk = 20
			dbSize.DiskUnits = "GiB"
			dbSize.PriceHourly = usdCost
			dbSize.CloudProvider = "aws"
			dbSize.Vcpus = vcpus
			dbSize.Region = UsWest1
			dbSize.PriceMonthly = usdCost * 730
			dbSize.Memory = memInt * 1024
			dbSize.MemoryUnits = "MiB"
			dbSize.Description = desc
			n = append(n, dbSize)
		}
	}

	err := hestia_compute_resources.InsertNodes(ctx, n)
	s.Require().NoError(err)
}

func TestAwsPricingClientTestSuite(t *testing.T) {
	suite.Run(t, new(AwsPricingClientTestSuite))
}
