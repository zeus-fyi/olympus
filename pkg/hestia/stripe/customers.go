package hestia_stripe

import (
	"context"
	"fmt"

	"github.com/stripe/stripe-go/v74"
)

func CreateCustomer(ctx context.Context, name, email string) error {
	params := &stripe.CustomerParams{
		Address:             nil,
		Balance:             nil,
		CashBalance:         nil,
		Coupon:              nil,
		DefaultSource:       nil,
		Description:         nil,
		Email:               stripe.String(email),
		InvoicePrefix:       nil,
		InvoiceSettings:     nil,
		Name:                stripe.String(name),
		NextInvoiceSequence: nil,
		PaymentMethod:       nil,
		Phone:               nil,
		PreferredLocales:    nil,
		PromotionCode:       nil,
		Shipping:            nil,
		Source:              nil,
		Tax:                 nil,
		TaxExempt:           nil,
		TaxIDData:           nil,
		TestClock:           nil,
		Validate:            nil,
	}
	//c, _ := customer.New(params)
	fmt.Println(params)
	return nil
}
