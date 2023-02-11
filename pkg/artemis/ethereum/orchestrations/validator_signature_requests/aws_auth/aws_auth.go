package artemis_hydra_orchestrations_auth

import (
	"context"
	"github.com/rs/zerolog/log"
	aws_secrets "github.com/zeus-fyi/zeus/pkg/aegis/aws"
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

func GetServiceRoutesAuths(ctx context.Context, p any) {

}
