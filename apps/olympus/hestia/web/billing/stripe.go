package hestia_billing

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/setupintent"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/keys"
)

type StripeBillingRequest struct {
}

func StripeBillingRequestHandler(c echo.Context) error {
	request := new(StripeBillingRequest)
	if err := c.Bind(request); err != nil {
		return err
	}
	return request.GetCustomerID(c)
}

type CheckoutData struct {
	ClientSecret string `json:"clientSecret"`
}

func (s *StripeBillingRequest) GetCustomerID(c echo.Context) error {
	ctx := context.Background()
	ou := c.Get("orgUser").(org_users.OrgUser)
	k := read_keys.OrgUserKey{
		OrgID: ou.OrgID,
		Key: keys.Key{
			UsersKeys: autogen_bases.UsersKeys{
				UserID: ou.UserID,
			},
		},
	}
	cID, err := k.GetOrCreateCustomerStripeID(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("GetOrCreateCustomerStripeID error")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	// &stripe.SetupIntentPaymentMethodDataParams{},
	params := &stripe.SetupIntentParams{
		Params: stripe.Params{
			Context:        nil,
			Expand:         nil,
			Extra:          nil,
			Headers:        nil,
			IdempotencyKey: nil,
			Metadata:       nil,
			StripeAccount:  nil,
		},
		AutomaticPaymentMethods: nil,
		ClientSecret:            nil,
		Customer:                stripe.String(cID),
		PaymentMethod:           stripe.String("card"),
		PaymentMethodData:       nil,
		PaymentMethodOptions:    nil,
		PaymentMethodTypes:      []*string{stripe.String("card")},
		ReturnURL:               nil,
		SingleUse:               nil,
		Usage:                   nil,
	}
	result, err := setupintent.New(params)
	if err != nil {
		log.Ctx(ctx).Err(err).Interface("ou", ou).Msg("setupintent.New error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	resp := CheckoutData{result.ClientSecret}
	return c.JSON(http.StatusOK, resp)
}
