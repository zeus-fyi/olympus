package hestia_stripe

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
)

func CreateCustomer(ctx context.Context, userID int, firstName, lastName, email string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(firstName + " " + lastName),
	}
	c, err := customer.New(params)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("CreateCustomer")
		return nil, err
	}
	k := create_keys.Key{}
	k.PublicKeyVerified = true
	k.PublicKeyName = "stripeCustomerID"
	k.PublicKey = c.ID
	k.PublicKeyTypeID = keys.StripeCustomerID
	k.UserID = userID
	err = k.InsertUserKey(ctx)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("CreateCustomer")
		_, derr := customer.Del(c.ID, nil)
		if derr != nil {
			log.Ctx(ctx).Err(derr).Msg("CreateCustomer, Delete Customer Cleanup")
		}
		return nil, err
	}
	return c, nil
}
