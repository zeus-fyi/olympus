package memfs

import (
	"errors"
	"io/fs"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (m *MemFS) ReadFileOutPath(p *filepaths.Path) ([]byte, error) {
	var b []byte
	if p == nil {
		return b, errors.New("need to include a path")
	}
	b, err := fs.ReadFile(m, p.FileOutPath())
	if err != nil {
		log.Err(err).Msgf("ReadFileOutPath %s", p.FileOutPath())
		return b, err
	}
	return b, nil
}

func (m *MemFS) ReadFileInPath(p *filepaths.Path) ([]byte, error) {
	var b []byte
	if p == nil {
		return b, errors.New("need to include a path")
	}
	b, err := fs.ReadFile(m, p.FileInPath())
	if err != nil {
		log.Err(err).Msgf("ReadFileInPath %s", p.FileInPath())
		return b, err
	}
	return b, nil
}
