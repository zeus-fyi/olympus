package artemis_hydra_orchestrations_auth

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/rs/zerolog/log"
	aws_secrets "github.com/zeus-fyi/zeus/pkg/aegis/aws"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

var HydraSecretManagerAuthAWS aws_secrets.SecretsManagerAuthAWS

func InitHydraSecretManagerAuthAWS(ctx context.Context, awsAuth aws_secrets.AuthAWS) {
	awsSm, err := aws_secrets.InitSecretsManager(ctx, awsAuth)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("Hydra: InitHestiaSecretManagerAuthAWS: error initializing aws secrets manager")
		panic(err)
	}
	HydraSecretManagerAuthAWS = awsSm
}

func GetServiceRoutesAuths(ctx context.Context, si aws_secrets.SecretInfo) (hestia_req_types.ServiceRequestWrapper, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(si.Name),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	srw := hestia_req_types.ServiceRequestWrapper{}
	result, err := HydraSecretManagerAuthAWS.GetSecretValue(ctx, input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Ctx(ctx).Err(err)
		return srw, err
	}
	// Decrypts secret using the associated KMS key.
	err = json.Unmarshal(result.SecretBinary, &srw)
	if err != nil {
		log.Ctx(ctx).Err(err)
		return srw, err
	}

	return srw, nil
}
