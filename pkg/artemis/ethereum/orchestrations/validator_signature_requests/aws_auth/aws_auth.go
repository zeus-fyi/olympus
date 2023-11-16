package artemis_hydra_orchestrations_aws_auth

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/rs/zerolog/log"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	aegis_aws_secretmanager "github.com/zeus-fyi/zeus/pkg/aegis/aws/secretmanager"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

var HydraSecretManagerAuthAWS aegis_aws_secretmanager.SecretsManagerAuthAWS

func InitHydraSecretManagerAuthAWS(ctx context.Context, awsAuth aegis_aws_auth.AuthAWS) {
	awsSm, err := aegis_aws_secretmanager.InitSecretsManager(ctx, awsAuth)
	if err != nil {
		log.Err(err).Msg("Hydra: InitHestiaSecretManagerAuthAWS: error initializing aws secrets manager")
		panic(err)
	}
	HydraSecretManagerAuthAWS = awsSm
	log.Info().Msg("InitHydraSecretManagerAuthAWS: initialized")
}

func GetServiceRoutesAuths(ctx context.Context, si aegis_aws_secretmanager.SecretInfo) (hestia_req_types.ServiceRequestWrapper, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(si.Name),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	srw := hestia_req_types.ServiceRequestWrapper{}
	result, err := HydraSecretManagerAuthAWS.GetSecretValue(ctx, input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Err(err).Msg("Hydra: GetServiceRoutesAuths: error getting secret value")
		return srw, err
	}
	// Decrypts secret using the associated KMS key.
	err = json.Unmarshal(result.SecretBinary, &srw)
	if err != nil {
		log.Err(err).Msg("Hydra: GetServiceRoutesAuths: error unmarshaling secret value")
		return srw, err
	}
	return srw, nil
}

func GetOrgSecret(ctx context.Context, name string) ([]byte, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(name),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	result, err := HydraSecretManagerAuthAWS.GetSecretValue(ctx, input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Err(err).Msg("Hydra: GetServiceRoutesAuths: error getting secret value")
		return nil, err
	}
	return result.SecretBinary, nil
}
