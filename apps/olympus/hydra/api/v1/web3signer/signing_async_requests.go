package hydra_eth2_web3signer

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	artemis_validator_signature_service_routing "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/signature_routing"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	bls_serverless_signing "github.com/zeus-fyi/zeus/pkg/aegis/aws/serverless_signing"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
)

type Resty struct {
	*resty.Client
}

// TODO add additional methods of usage besides aws lambda

func RequestValidatorSignaturesAsync(ctx context.Context, sigRequests aegis_inmemdbs.EthereumBLSKeySignatureRequests, pubkeyToUUID map[string]string) error {
	gm := artemis_validator_signature_service_routing.GroupSigRequestsByGroupName(ctx, sigRequests)
	for groupName, signReqs := range gm {
		go func(groupName string, signReqs aegis_inmemdbs.EthereumBLSKeySignatureRequests) {
			auth, err := artemis_validator_signature_service_routing.GetGroupAuthFromInMemFS(ctx, groupName)
			if err != nil {
				log.Ctx(ctx).Err(err).Msg("Failed to get group auth")
				return
			}
			sr := bls_serverless_signing.SignatureRequests{
				SecretName:        auth.SecretName,
				SignatureRequests: aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: signReqs.Map},
			}
			sigResponses := aegis_inmemdbs.EthereumBLSKeySignatureResponses{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureResponse)}

			cfg := aegis_aws_auth.AuthAWS{
				Region:    "us-west-1",
				AccessKey: auth.AccessKey,
				SecretKey: auth.SecretKey,
			}
			r := Resty{}
			r.Client = resty.New()
			r.SetTimeout(2 * time.Second)
			r.SetRetryCount(3)
			minDuration := 10 * time.Millisecond
			maxDuration := 100 * time.Millisecond
			jitter := time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration
			r.SetRetryWaitTime(jitter)
			r.SetBaseURL(auth.AuthLamdbaAWS.ServiceURL)
			reqAuth, err := cfg.CreateV4AuthPOSTReq(ctx, "lambda", auth.AuthLamdbaAWS.ServiceURL, sr)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("failed to get service routes auths for lambda iam auth")
				return
			}
			resp, err := r.R().
				SetHeaderMultiValues(reqAuth.Header).
				SetResult(&sigResponses).
				SetBody(sr).Post("/")

			// TODO, notify on errors, track these metrics & latency
			if err != nil {
				log.Ctx(ctx).Err(err).Msg("Failed to get response")
				return
			}
			if resp.StatusCode() != 200 {
				err = errors.New("non-200 status code")
				log.Ctx(ctx).Err(err).Msg("Failed to get 200 status code")
				return
			}

			if len(sigRequests.Map) < len(sigRequests.Map) {
				log.Ctx(ctx).Warn().Msg("Not all signatures were returned")
				log.Ctx(ctx).Info().Interface("sigRequests", sigRequests).Interface("sigResponses", sigResponses).Msg("Not all signatures were returned")
			}
			for pubkey, sigRespWrapper := range sigResponses.Map {
				uuid := pubkeyToUUID[pubkey]
				SignatureResponsesCache.Set(uuid, sigRespWrapper.Signature, cache.DefaultExpiration)
			}
		}(groupName, signReqs)
	}
	return nil
}
