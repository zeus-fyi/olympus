package filepaths

import (
	"os"

	"github.com/rs/zerolog/log"
)

func (p *Path) OpenFileInPath() (*os.File, error) {
	f, err := os.Open(p.FileInPath())
	if err != nil {
		log.Err(err).Msg("OpenFileInPath")
		return nil, err
	}
	return f, err
}
