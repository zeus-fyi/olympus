package snapshot_init

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/configs"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
)

func SetConfigByEnv(ctx context.Context, env string) {
	switch env {
	case "production":
		log.Info().Msg("Downloader: production auth procedure starting")
		authCfg := auth_startup.NewDefaultAuthClient(ctx, authKeysCfg)
		_, sw := auth_startup.RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx, authCfg)
		cfg.PGConnStr = sw.PostgresAuth
	case "production-local":
		tc := configs.InitLocalTestConfigs()
		authKeysCfg = tc.ProdLocalAuthKeysCfg
		cfg.PGConnStr = tc.ProdLocalDbPgconn
	case "local":
		tc := configs.InitLocalTestConfigs()
		cfg.PGConnStr = tc.LocalDbPgconn
	}
}
