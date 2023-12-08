package hestia_billing

import (
	"context"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	hestia_stripe "github.com/zeus-fyi/olympus/pkg/hestia/stripe"
)

var BillingCache = cache.New(time.Hour, cache.DefaultExpiration)

func CheckBillingCache(ctx context.Context, userID int) bool {
	billingExists, ok := BillingCache.Get(fmt.Sprintf("%d", userID))
	if ok {
		b, bok := billingExists.(bool)
		if bok && b {
			log.Info().Interface("userID", userID).Interface("billingExists", billingExists).Msg("found in cache")
			return b
		}
	}
	isBillingSetup, err := hestia_stripe.DoesUserHaveBillingMethod(ctx, userID)
	if err != nil {
		log.Error().Err(err).Interface("userID", userID).Msg("failed to check if user has billing method")
		return false
	}
	BillingCache.Set(fmt.Sprintf("%d", userID), isBillingSetup, 1*time.Hour)
	return isBillingSetup
}
