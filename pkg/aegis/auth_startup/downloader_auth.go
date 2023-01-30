package auth_startup

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

func InitDownloaderAuth(ctx context.Context, inMemSecrets memfs.MemFS, secrets SecretsWrapper) {
	log.Info().Msg("Downloader: InitDownloaderAuth starting")
	secrets.PostgresAuth = secrets.ReadSecret(ctx, inMemSecrets, pgSecret)
	log.Info().Msg("Downloader: InitDownloaderAuth done")
	return
}
