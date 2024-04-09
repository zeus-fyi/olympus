package hestia_stripe

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/paymentmethod"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateCustomer(ctx context.Context, userID int, firstName, lastName, email string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(firstName + " " + lastName),
	}
	c, err := customer.New(params)
	if err != nil {
		log.Err(err).Msg("CreateCustomer")
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
		log.Err(err).Msg("CreateCustomer")
		_, derr := customer.Del(c.ID, nil)
		if derr != nil {
			log.Err(derr).Msg("CreateCustomer, Delete Customer Cleanup")
		}
		return nil, err
	}
	return c, nil
}

const (
	InternalUserID        = 7138958574876245565
	InternalUserID1       = 7138958574876245567
	InternalUserID2       = 1699642242976434000
	InternalSharedUserID3 = 1710298581127603000

	CustomerDemoOrgID = 1712629107490733000
)

func DoesUserHaveBillingMethod(ctx context.Context, userID int) (bool, error) {
	if userID == CustomerDemoOrgID || userID == InternalUserID || userID == InternalUserID1 || userID == InternalUserID2 || userID == InternalSharedUserID3 {
		return true, nil
	}
	cID, err := QueryGetCustomerStripeID(ctx, userID)
	if err != nil {
		log.Err(err).Interface("u", userID).Msg("DoesUserHaveBillingMethod error")
		return false, nil
	}
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(cID),
		Type:     stripe.String("card"),
	}
	i := paymentmethod.List(params)
	for i.Next() {
		return true, nil
	}
	return false, nil
}

func QueryGetCustomerStripeID(ctx context.Context, userID int) (string, error) {
	var q sql_query_templates.QueryParams
	query := fmt.Sprintf(`
	SELECT usk.public_key
	FROM users_keys usk
	INNER JOIN key_types kt ON kt.key_type_id = usk.public_key_type_id
	INNER JOIN org_users ou ON ou.user_id = usk.user_id
	INNER JOIN users u ON u.user_id = ou.user_id
	WHERE u.user_id = $1 AND usk.public_key_type_id = $2
	`)
	q.RawQuery = query
	pubkey := ""
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, userID, keys.StripeCustomerID).Scan(&pubkey)
	if err != nil {
		return "", err
	}
	return pubkey, nil
}
