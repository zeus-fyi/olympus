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
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()
	gm := artemis_validator_signature_service_routing.GroupSigRequestsByGroupName(ctx, sigRequests)
	for groupName, signReqs := range gm {
		go func(groupName string, signReqs aegis_inmemdbs.EthereumBLSKeySignatureRequests) {
			auth, err := artemis_validator_signature_service_routing.GetGroupAuthFromInMemFS(ctx, groupName)
			if err != nil {
				log.Err(err).Msg("Failed to get group auth")
				return
			}
			sr := bls_serverless_signing.SignatureRequests{
				SecretName:        auth.SecretName,
				SignatureRequests: aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: signReqs.Map},
			}
			cfg := aegis_aws_auth.AuthAWS{
				Region:    "us-west-1",
				AccessKey: auth.AccessKey,
				SecretKey: auth.SecretKey,
			}

			ch := make(chan aegis_inmemdbs.EthereumBLSKeySignatureResponses, 1)
			for i := 0; i < 3; i++ {
				go func(i int) {
					switch i {
					case 1:
						minDuration := 150 * time.Millisecond
						maxDuration := 200 * time.Millisecond
						jitter := time.Duration(i) * (time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration)
						time.Sleep(jitter)
					case 2:
						minDuration := 200 * time.Millisecond
						maxDuration := 250 * time.Millisecond
						jitter := time.Duration(i) * (time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration)
						time.Sleep(jitter)
					case 3:
						minDuration := 250 * time.Millisecond
						maxDuration := 300 * time.Millisecond
						jitter := time.Duration(i) * (time.Duration(rand.Int63n(int64(maxDuration-minDuration))) + minDuration)
						time.Sleep(jitter)
					default:
					}
					if len(ch) == cap(ch) {
						return
					}
					sigResponses := aegis_inmemdbs.EthereumBLSKeySignatureResponses{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureResponse)}
					r := Resty{}
					r.Client = resty.New()
					r.SetBaseURL(auth.AuthLamdbaAWS.ServiceURL)
					reqAuth, rerr := cfg.CreateV4AuthPOSTReq(ctx, "lambda", auth.AuthLamdbaAWS.ServiceURL, sr)
					if rerr != nil {
						log.Error().Err(rerr).Msg("failed to get service routes auths for lambda iam auth")
						return
					}
					resp, reerr := r.R().
						SetHeaderMultiValues(reqAuth.Header).
						SetResult(&sigResponses).
						SetBody(sr).Post("/")
					if reerr != nil {
						log.Err(reerr).Msg("Failed to get response")
						return
					}
					if resp.StatusCode() != 200 {
						err = errors.New("non-200 status code")
						log.Err(reerr).Msg("Failed to get 200 status code")
						return
					}
					ch <- sigResponses
				}(i)
			}
			timeout := time.After(6 * time.Second)
			select {
			case sigResponses := <-ch:
				if len(sigRequests.Map) < len(sigRequests.Map) {
					log.Warn().Msg("Not all signatures were returned")
					log.Info().Interface("sigRequests", sigRequests).Interface("sigResponses", sigResponses).Msg("Not all signatures were returned")
				}
				for pubkey, sigRespWrapper := range sigResponses.Map {
					uuid := pubkeyToUUID[pubkey]
					SignatureResponsesCache.Set(uuid, sigRespWrapper.Signature, cache.DefaultExpiration)
				}
				return
			case <-timeout:
				log.Error().Msg("Timed out")
			}
		}(groupName, signReqs)
	}
	return nil
}
