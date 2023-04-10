package hestia_billing

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/setupintent"
)

type StripeBillingRequest struct {
}

func StripeBillingRequestHandler(c echo.Context) error {
	request := new(StripeBillingRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.Dothing(c)
}

type CheckoutData struct {
	ClientSecret string `json:"client_secret"`
}

func (s *StripeBillingRequest) Dothing(c echo.Context) error {

	// GET customer ID from session
	params := &stripe.SetupIntentParams{
		Customer:           stripe.String("{{CUSTOMER_ID}}"),
		PaymentMethodTypes: []*string{stripe.String("card"), stripe.String("link")},
	}
	result, err := setupintent.New(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp := CheckoutData{result.ClientSecret}
	return c.JSON(http.StatusOK, resp)
}
