package hestia_stripe

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/product"
)

func CreateMeteredProductWithPricing(ctx context.Context, productName, productDescription string, pricingHourly, pricingMonthly int64) error {
	params := &stripe.ProductParams{
		Name:        stripe.String(productName),
		Description: stripe.String(productDescription),
	}
	prod, err := product.New(params)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("hestia_stripe.CreateMeteredProductPricing Error creating product")
		return nil
	}
	fmt.Printf("Product created: %v\n", prod.ID)
	pricingParamsHourly := &stripe.PriceParams{
		Nickname:   stripe.String(fmt.Sprintf("Per %s pricing", "hour")),
		UnitAmount: stripe.Int64(pricingHourly),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		Recurring: &stripe.PriceRecurringParams{
			Interval:  stripe.String(string(stripe.PriceRecurringIntervalMonth)),
			UsageType: stripe.String(string(stripe.PriceRecurringUsageTypeMetered)),
		},
		Product: stripe.String(prod.ID),
	}
	pricingParamsDaily := &stripe.PriceParams{
		Nickname:   stripe.String(fmt.Sprintf("Per %s pricing", "24 hours")),
		UnitAmount: stripe.Int64(pricingHourly * 24),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		Recurring: &stripe.PriceRecurringParams{
			Interval:  stripe.String(string(stripe.PriceRecurringIntervalMonth)),
			UsageType: stripe.String(string(stripe.PriceRecurringUsageTypeMetered)),
		},
		Product: stripe.String(prod.ID),
	}
	pricingParamsMonthly := &stripe.PriceParams{
		Nickname:   stripe.String(fmt.Sprintf("Per %s pricing", "month")),
		UnitAmount: stripe.Int64(pricingMonthly),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		Recurring: &stripe.PriceRecurringParams{
			Interval:  stripe.String(string(stripe.PriceRecurringIntervalMonth)),
			UsageType: stripe.String(string(stripe.PriceRecurringUsageTypeMetered)),
		},
		Product: stripe.String(prod.ID),
	}
	result, err := price.New(pricingParamsHourly)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("hestia_stripe.CreateMeteredProductPricing")
		return nil
	}
	fmt.Printf("Price created: %v\n", result.ID)
	result, err = price.New(pricingParamsDaily)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("hestia_stripe.CreateMeteredProductPricing")
		return nil
	}
	fmt.Printf("Price created: %v\n", result.ID)
	result, err = price.New(pricingParamsMonthly)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("hestia_stripe.CreateMeteredProductPricing")
		return nil
	}
	fmt.Printf("Price created: %v\n", result.ID)
	return nil
}

func ArchiveProducts(ctx context.Context) error {
	// Retrieve a list of all products
	productListParams := &stripe.ProductListParams{
		ListParams: stripe.ListParams{
			// Set a high limit to retrieve all products
			Limit: stripe.Int64(100),
		},
	}
	products := product.List(productListParams)
	// Loop through the list and archive each product
	for products.Next() {
		productParams := &stripe.ProductParams{
			Active: stripe.Bool(false),
		}
		_, err := product.Update(products.Product().ID, productParams)
		if err != nil {
			return err
		}
	}
	if err := products.Err(); err != nil {
		return err
	}
	fmt.Println("All products archived successfully.")
	return nil
}
