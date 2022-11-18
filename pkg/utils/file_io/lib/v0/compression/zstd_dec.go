package compression

import (
	"errors"
	"os"

	"github.com/klauspost/compress/zstd"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (c *Compression) ZstdDecompress(p *filepaths.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	r, err := os.Open(p.FileInPath())
	if err != nil {
		log.Err(err).Msg("Compression: ZstdDecompress, os.Open(p.FileInPath())")
		return err
	}
	defer r.Close()
	zstdReader, err := zstd.NewReader(r)
	if err != nil {
		log.Err(err).Msg("Compression: ZstdDecompress, zstd.NewReader(r)")
		return err
	}
	defer zstdReader.Close()
	return tarReader(p, zstdReader)
}
