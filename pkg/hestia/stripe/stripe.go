package hestia_stripe

import "github.com/stripe/stripe-go/v74"

func InitStripe(apiKey string) {
	stripe.Key = apiKey
}

//func CreateIntent(ctx context.Context, customerID, checkoutSessionID string) (*stripe.CheckoutSession, error) {
//	params := &stripe.CheckoutSessionParams{
//		PaymentMethodTypes: stripe.StringSlice([]string{
//			"card",
//		}),
//		Mode:       stripe.String(string(stripe.CheckoutSessionModeSetup)),
//		Customer:   stripe.String(fmt.Sprintf("{{%s}}", customerID)),
//		SuccessURL: stripe.String(fmt.Sprintf("https://cloud.zeus.fyi/success?session_id={%s}", checkoutSessionID)),
//		CancelURL:  stripe.String("https://cloud.zeus.fyi/cancel"),
//	}
//	cs, err := session.New(params)
//	if err != nil {
//		log.Ctx(ctx).Error().Err(err).Msg("hestia_stripe.CreateIntent")
//		return nil, err
//	}
//	return cs, nil
//}
