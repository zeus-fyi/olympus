package filepaths

import (
	"os"

	"github.com/rs/zerolog/log"
)

func (p *Path) ReadFileInPath() ([]byte, error) {
	byteArray, err := os.ReadFile(p.FileInPath())
	if err != nil {
		log.Err(err).Msg("ReadFileInPath")
		return []byte{}, err
	}
	return byteArray, err
}

func (p *Path) ReadFileOutPath() ([]byte, error) {
	byteArray, err := os.ReadFile(p.FileOutPath())
	if err != nil {
		log.Err(err).Msg("FileOutPath")
		return []byte{}, err
	}
	return byteArray, err
}
