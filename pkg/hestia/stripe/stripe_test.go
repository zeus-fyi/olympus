package hestia_stripe

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_digitalocean "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type StripeTestSuite struct {
	do hestia_digitalocean.DigitalOcean
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (s *StripeTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.do = hestia_digitalocean.InitDoClient(ctx, s.Tc.DigitalOceanAPIKey)
	s.Require().NotNil(s.do.Client)
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	InitStripe(s.Tc.StripeTestSecretAPIKey)
}

func (s *StripeTestSuite) TestArchiveProducts() {
	err := ArchiveProducts(ctx)
	s.Require().NoError(err)
}

func (s *StripeTestSuite) TestCreateProduct() {
	sizes, err := s.do.GetSizes(ctx)
	s.Require().NoError(err)
	s.Require().NotEmpty(sizes)
	for _, size := range sizes {
		hourlyPricing := math.Ceil(size.PriceHourly * 100 * 1.1)
		monthlyPricing := size.PriceMonthly * 100 * 1.1
		fmt.Println(hourlyPricing, monthlyPricing, size.Slug)
		description := fmt.Sprintf("DigitalOcean %s", size.Description)
		err = CreateMeteredProductWithPricing(ctx, size.Slug, description, int64(hourlyPricing), int64(monthlyPricing))
		s.Require().NoError(err)
	}
	productName := "DigitalOcean per 100Gi Block Storage SSD"
	hourlyPricing := math.Ceil(0.015 * 100 * 1.1)
	dailyPricing := hourlyPricing * 24 * 30
	monthlyPricing := 10 * 100 * 1.1
	fmt.Println(hourlyPricing, dailyPricing, monthlyPricing, productName)
	err = CreateMeteredProductWithPricing(ctx, productName, productName, int64(hourlyPricing), int64(monthlyPricing))
	s.Require().NoError(err)
}

func TestStripeTestSuite(t *testing.T) {
	suite.Run(t, new(StripeTestSuite))
}
