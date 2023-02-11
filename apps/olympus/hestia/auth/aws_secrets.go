package hestia_aws_secrets_auth

import (
	"context"
	"github.com/rs/zerolog/log"
	aws_secrets "github.com/zeus-fyi/zeus/pkg/aegis/aws"
)

var HestiaSecretManagerAuthAWS = aws_secrets.SecretsManagerAuthAWS{}

func InitHestiaSecretManagerAuthAWS(ctx context.Context, awsAuth aws_secrets.AuthAWS) {
	awsSm, err := aws_secrets.InitSecretsManager(ctx, awsAuth)
	if err != nil {
		log.Ctx(ctx).Panic().Err(err).Msg("Hestia: InitHestiaSecretManagerAuthAWS: error initializing aws secrets manager")
		panic(err)
	}
	HestiaSecretManagerAuthAWS = awsSm
}
