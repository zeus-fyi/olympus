package compression

import (
	"compress/gzip"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

// GzipDecompress takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func (c *Compression) GzipDecompress(p *filepaths.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}

	r, err := os.Open(p.FileInPath())
	if err != nil {
		return err
	}
	defer r.Close()
	gzr, err := gzip.NewReader(r)
	if err != nil {
		log.Err(err).Msg("Compression: GzipDecompress, gzip.NewReader(r))")
		return err
	}
	defer gzr.Close()
	return tarReader(p, gzr)
}
