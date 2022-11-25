package compression

import (
	"archive/tar"
	"errors"
	"io/fs"
	"os"

	"github.com/pierrec/lz4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (c *Compression) Lz4CompressDir(p *filepaths.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}

	p.FnOut = p.FnIn + ".tar.lz4"
	out, err := os.Create(p.FileOutPath())
	if err != nil {
		log.Err(err).Msg("Compression: Lz4CompressDir, os.Create(p.FileOutPath())")
		return err
	}
	defer out.Close()

	enc := lz4.NewWriter(out)
	defer enc.Close()
	tw := tar.NewWriter(enc)
	defer tw.Close()

	fileSystem := os.DirFS(p.DirIn)
	err = fs.WalkDir(fileSystem, ".", func(filename string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Err(err).Msgf("Compression: fs.WalkDir at filename %s", filename)
			return err
		}
		if !d.IsDir() {
			zerr := addToArchive(p, tw, filename)
			if zerr != nil {
				log.Err(zerr).Msgf("Compression: addToArchive at filename %s", filename)
				return zerr
			}
		}
		return nil
	})

	p.DirIn = p.DirOut
	p.FnIn = p.FnOut
	return err
}
