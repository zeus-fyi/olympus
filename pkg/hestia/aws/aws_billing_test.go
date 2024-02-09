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
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

type AwsPricingClientTestSuite struct {
	test_suites_base.TestSuite
	pc AwsPricing
}

func (s *AwsPricingClientTestSuite) SetupTest() {
	s.InitLocalConfigs()
	eksCreds := aegis_aws_auth.AuthAWS{
		Region:    "us-east-1",
		AccessKey: s.Tc.AwsAccessKeyEks,
		SecretKey: s.Tc.AwsSecretKeyEks,
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

/*
gen 3

i3.4xlarge	$1.376	16	122 GiB	2 x 1900 NVMe SSD	Up to 10 Gigabit
i3.8xlarge	$2.752	32	244 GiB	4 x 1900 NVMe SSD	10 Gigabit

gen 4
i4i.4xlarge	$1.514	16	128 GiB	1 x 3750 NVMe SSD	Up to 25 Gigabit
i4i.8xlarge	$3.027	32	256 GiB	2 x 3750 NVMe SSD	18750 Megabit
*/

func (s *AwsPricingClientTestSuite) TestGetEC2Product() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	instanceTypes := []string{
		//"t3.nano",
		//"t3.micro",
		//"t3.small",
		//"t3.medium",
		//"t3.large",
		//"t3.xlarge",
		//"t3.2xlarge",
		//"i3.4xlarge",
		//"i3.8xlarge",
		//"i4i.4xlarge",
		//"i4i.8xlarge",
		//"c3.4xlarge",
		//"c1.medium",
		//"c1.xlarge",
		//"c3.2xlarge",
		//"c3.4xlarge",
		//"c3.8xlarge",
		//"c3.large",
		//"c3.xlarge",
		//"c4.large",
		//"c4.xlarge",
		//"c5.12xlarge",
		//"c5.18xlarge",
		//"c5.24xlarge",
		//"c5.2xlarge",
		//"c5.4xlarge",
		"m7g.medium",
		"m7g.large",
		"m7g.xlarge",
		"m7g.2xlarge",
		"m7g.4xlarge",
		"m7g.8xlarge",
	}
	region := "us-east-1"

	diskType := "ssd"

	n := hestia_autogen_bases.NodesSlice{}
	for _, instanceType := range instanceTypes {
		prices, err := s.pc.GetEC2Product(ctx, region, instanceType)
		s.Require().NoError(err)
		fmt.Printf("%s\n", instanceType)
		for _, price := range prices {
			desc := price.GetDescription()
			if !strings.Contains(desc, fmt.Sprintf("per On Demand Linux %s Instance Hour", instanceType)) {
				continue
			}
			fmt.Println(desc)
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
			dbSize.DiskUnits = "GB"
			dbSize.DiskType = diskType
			dbSize.PriceHourly = usdCost
			dbSize.CloudProvider = "aws"
			dbSize.Vcpus = vcpus
			dbSize.Region = region
			dbSize.PriceMonthly = usdCost * 730
			dbSize.Memory = memInt * 1024
			dbSize.MemoryUnits = "MiB"
			dbSize.Description = desc
			n = append(n, dbSize)
		}
	}

	s.Require().NotEmpty(n)
	//
	err := hestia_compute_resources.InsertNodes(ctx, n)
	s.Require().NoError(err)
}

func TestAwsPricingClientTestSuite(t *testing.T) {
	suite.Run(t, new(AwsPricingClientTestSuite))
}
