package eth_validator_signature_requests

import (
	"context"
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	artemis_validator_signature_service_routing "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/signature_routing"
	eth_validators_service_requests "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validators_service_requests"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	bls_serverless_signing "github.com/zeus-fyi/zeus/pkg/aegis/aws/serverless_signing"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
)

const (
	waitForTxRxTimeout    = 15 * time.Minute
	submitSignedTxTimeout = 5 * time.Minute
)

type ArtemisEthereumValidatorSignatureRequestActivities struct {
}

func NewArtemisEthereumValidatorSignatureRequestActivities() ArtemisEthereumValidatorSignatureRequestActivities {
	return ArtemisEthereumValidatorSignatureRequestActivities{}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

type Eth2SignResponse struct {
	Signature string `json:"signature"`
}

func (d *ArtemisEthereumValidatorSignatureRequestActivities) GetActivities() ActivitiesSlice {
	return []interface{}{d.RequestValidatorSignatures, d.SendHeartbeat}
}

type Resty struct {
	*resty.Client
}

// TODO add signature auth, todo make each group an async activity
func (d *ArtemisEthereumValidatorSignatureRequestActivities) RequestValidatorSignatures(ctx context.Context, sigRequests aegis_inmemdbs.EthereumBLSKeySignatureRequests) (aegis_inmemdbs.EthereumBLSKeySignatureResponses, error) {
	sigResponses := aegis_inmemdbs.EthereumBLSKeySignatureResponses{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureResponse)}
	gm := artemis_validator_signature_service_routing.GroupSigRequestsByGroupName(ctx, sigRequests)
	r := Resty{}
	r.Client = resty.New()
	for groupName, signReqs := range gm {
		auth, err := artemis_validator_signature_service_routing.GetGroupAuthFromInMemFS(ctx, groupName)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to get group auth")
			return sigResponses, err
		}
		sr := bls_serverless_signing.SignatureRequests{
			SecretName:        auth.SecretName,
			SignatureRequests: aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: signReqs.Map},
		}
		respJson := aegis_inmemdbs.EthereumBLSKeySignatureResponses{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureResponse)}
		_, err = r.R().
			SetResult(&respJson).
			SetBody(sr).
			Post(auth.AuthLamdbaAWS.ServiceURL)
		if err != nil {
			log.Ctx(ctx).Error().Err(err)
			return sigResponses, err
		}
		for k, v := range respJson.Map {
			sigResponses.Map[k] = v
		}
	}

	if len(sigRequests.Map) < len(sigRequests.Map) {
		log.Ctx(ctx).Warn().Msg("Not all signatures were returned")
		log.Ctx(ctx).Info().Interface("sigRequests", sigRequests).Interface("sigResponses", sigResponses).Msg("Not all signatures were returned")
		return sigResponses, nil
	}

	return sigResponses, nil
}

func (d *ArtemisEthereumValidatorSignatureRequestActivities) SendHeartbeat(ctx context.Context) ([]string, error) {
	log.Ctx(ctx).Info().Msg("sending heartbeat message")
	var restore []string
	for {
		ql := HeartbeatQueue.Size()
		if ql == 0 {
			break
		}
		groupName, qOk := HeartbeatQueue.Dequeue()
		if !qOk {
			continue
		}
		restore = append(restore, groupName)
		log.Ctx(ctx).Info().Interface("groupName", groupName).Msg("sending heartbeat message")
		go func(groupName string) {
			log.Ctx(ctx).Info().Interface("groupName", groupName).Msg("sending heartbeat message")
			auth, err := artemis_validator_signature_service_routing.GetGroupAuthFromInMemFS(ctx, groupName)
			if err != nil {
				log.Ctx(ctx).Err(err).Msg("Failed to get group auth")
				return
			}
			signReqs := aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}
			hexMessage, err := aegis_inmemdbs.RandomHex(10)
			if err != nil {
				log.Ctx(ctx).Err(err).Msg("failed to create random hex message")
				return
			}
			signReqs.Map["0x0000000"] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: hexMessage}
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
			r := eth_validators_service_requests.Resty{}
			r.Client = resty.New()
			r.SetBaseURL(auth.AuthLamdbaAWS.ServiceURL)
			r.SetTimeout(5 * time.Second)
			r.SetRetryCount(2)
			r.SetRetryWaitTime(20 * time.Millisecond)
			reqAuth, err := cfg.CreateV4AuthPOSTReq(ctx, "lambda", auth.AuthLamdbaAWS.ServiceURL, sr)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("failed to get service routes auths for lambda iam auth")
				return
			}
			log.Info().Interface("groupName", groupName).Msg("sending heartbeat")
			resp, err := r.R().
				SetHeaderMultiValues(reqAuth.Header).
				SetResult(&sigResponses).
				SetBody(sr).Post("/")
			if err != nil {
				log.Ctx(ctx).Err(err).Interface("groupName", groupName).Msg("failed to get response")
				return
			}
			if resp.StatusCode() != 200 {
				err = errors.New("non-200 status code")
				log.Ctx(ctx).Err(err).Interface("groupName", groupName).Msg("failed to get 200 status code")
				return
			} else {
				log.Ctx(ctx).Info().Interface("groupName", groupName).Msg("heartbeat OK")
			}
		}(groupName)
	}

	return restore, nil
}
