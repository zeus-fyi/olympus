package hydra_server

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
)

var temporalProdAuthConfig = temporal_auth.TemporalAuth{
	ClientCertPath:   "/etc/ssl/certs/ca.pem",
	ClientPEMKeyPath: "/etc/ssl/certs/ca.key",
	Namespace:        "production-artemis.ngb72",
	HostPort:         "production-artemis.ngb72.tmprl.cloud:7233",
}

func SetConfigByEnv(ctx context.Context, env string) {
	switch env {
	case "production":
	case "production-local":
	case "local":
	}

	log.Info().Msgf("Hydra %s temporal auth and init procedure starting", env)
	// TODO replace
	log.Info().Msgf("Hydra %s temporal auth and init procedure succeeded", env)
}
