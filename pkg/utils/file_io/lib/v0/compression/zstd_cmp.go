package compression

import (
	"archive/tar"
	"context"
	"errors"
	"io/fs"
	"os"

	"github.com/klauspost/compress/zstd"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (c *Compression) ZstCompressDir(p *filepaths.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}

	p.FnOut = p.FnIn + ".tar.zst"
	out, err := os.Create(p.FileOutPath())
	if err != nil {
		log.Err(err).Msg("Compression: ZstCompressDir, os.Create(p.FileOutPath())")
		return err
	}
	defer out.Close()

	enc, err := zstd.NewWriter(out)
	if err != nil {
		log.Err(err).Msg("Compression: ZstCompressDir, zstd.NewWriter(out)")
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
		if !d.IsDir() && filename != p.FnOut {
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

func (c *Compression) ZstCompressFile(ctx context.Context, p *filepaths.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}

	p.FnOut = p.FnIn + ".tar.zst"
	out, err := os.Create(p.FileOutPath())
	if err != nil {
		log.Err(err).Msg("Compression: ZstCompressDir, os.Create(p.FileOutPath())")
		return err
	}
	defer out.Close()

	enc, err := zstd.NewWriter(out)
	if err != nil {
		log.Err(err).Msg("Compression: ZstCompressDir, zstd.NewWriter(out)")
		return err
	}
	defer enc.Close()
	tw := tar.NewWriter(enc)
	defer tw.Close()

	// Compress the single file directly
	err = addToArchive(p, tw, p.FnIn)
	if err != nil {
		return err
	}

	// Update the path and file names to the output
	p.DirIn = p.DirOut
	p.FnIn = p.FnOut
	return nil
}
