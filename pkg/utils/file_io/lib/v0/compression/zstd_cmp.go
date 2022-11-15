package compression

import (
	"archive/tar"
	"errors"
	"io/fs"
	"os"

	"github.com/klauspost/compress/zstd"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (c *Compression) CreateTarZstdArchiveDir(p *structs.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}

	p.FnOut = p.FnIn + ".tar.zst"
	out, err := os.Create(p.FileOutPath())
	if err != nil {
		log.Err(err).Msg("Compression: CreateTarZstdArchiveDir, os.Create(p.FileOutPath())")
		return err
	}
	defer out.Close()

	enc, err := zstd.NewWriter(out)
	if err != nil {
		log.Err(err).Msg("Compression: CreateTarZstdArchiveDir, zstd.NewWriter(out)")
		return err
	}
	defer enc.Close()
	tw := tar.NewWriter(enc)
	defer tw.Close()

	fileSystem := os.DirFS(p.DirIn)
	err = fs.WalkDir(fileSystem, ".", func(filename string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			zerr := addToArchive(p, tw, filename)
			if zerr != nil {
				return zerr
			}
		}
		return nil
	})

	p.DirIn = p.DirOut
	p.FnIn = p.FnOut
	return err
}
