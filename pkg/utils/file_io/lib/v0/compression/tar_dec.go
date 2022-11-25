package compression

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (c *Compression) TarDecompress(p *filepaths.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	r, err := os.Open(p.FileInPath())
	if err != nil {
		log.Err(err).Msg("Compression: TarDecompress, os.Open(p.FileInPath())")
		return err
	}
	defer r.Close()
	return tarReader(p, r)
}
