package structs

import (
	"os"

	"github.com/rs/zerolog/log"
)

func (p *Path) OpenFileInPath() (*os.File, error) {
	f, err := os.OpenFile(p.FileInPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Err(err).Msg("OpenFileInPath")
		return nil, err
	}
	return f, err
}
