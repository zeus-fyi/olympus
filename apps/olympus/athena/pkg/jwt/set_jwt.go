package athena_jwt

import (
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func SetToken(p filepaths.Path, token string) error {
	var fileIO = file_io.FileIO{}

	p.FnOut = "jwt.hex"
	err := fileIO.CreateV2FileOut(p, []byte(token))
	if err != nil {
		log.Err(err).Msg("TokenRequestJWT.CreateV2FileOut(athena_server.DataDir, b)")
		return err
	}
	return err
}
