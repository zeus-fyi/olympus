package auth_startup

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

func NewDigitalOceanS3AuthClient(ctx context.Context, keysCfg auth_keys_config.AuthKeysCfg) s3base.S3Client {
	if len(keysCfg.SpacesPrivKey) <= 0 {
		log.Warn().Msg("no spaces priv key provided")
		misc.DelayedPanic(fmt.Errorf("no spaces priv key provided"))
	}
	if len(keysCfg.SpacesKey) <= 0 {
		log.Warn().Msg("no spaces key provided")
		misc.DelayedPanic(fmt.Errorf("no spaces key provided"))
	}
	s3BaseClient, err := s3base.NewConnS3ClientWithStaticCreds(ctx, keysCfg.SpacesKey, keysCfg.SpacesPrivKey)
	if err != nil {
		log.Fatal().Msg("NewDefaultAuthClient: NewConnS3ClientWithStaticCreds failed, shutting down the server")
		misc.DelayedPanic(err)
	}
	return s3BaseClient
}
