package compression

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/pierrec/lz4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

func (c *Compression) Lz4Decompress(p *filepaths.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	r, err := os.Open(p.FileInPath())
	if err != nil {
		log.Err(err).Msg("Compression: Lz4Decompress, os.Open(p.FileInPath())")
		return err
	}
	defer r.Close()
	lz4Reader := lz4.NewReader(r)
	if err != nil {
		log.Err(err).Msg("Compression: Lz4Decompress, lz4.NewReader(r)")
		return err
	}
	return tarReader(p, lz4Reader)
}

func (c *Compression) Lz4DecompressInMemFsFile(p *filepaths.Path, inMemFs memfs.MemFS) (memfs.MemFS, error) {
	if p == nil {
		return inMemFs, errors.New("need to include a path")
	}
	b, err := inMemFs.ReadFileInPath(p)
	if err != nil {
		log.Err(err).Msg("Lz4DecompressInMemFsFile: ")
		return inMemFs, err
	}
	o, err := decompress(b)
	err = inMemFs.MakeFileOut(p, o)
	p.DirIn = p.DirOut
	p.FnIn = p.FnOut
	if err != nil {
		log.Err(err).Msg("Lz4DecompressInMemFsFile: ")
		return inMemFs, err
	}
	return inMemFs, err
}

func compress(toCompress []byte) ([]byte, int, error) {
	toCompressLen := len(toCompress)
	compressed := make([]byte, toCompressLen)

	//compress
	l, err := lz4.CompressBlock(toCompress, compressed, nil)
	if err != nil {
		panic(err)
	}
	return toCompress[:l], toCompressLen, nil
}

func decompress(in []byte) ([]byte, error) {
	r := bytes.NewReader(in)
	w := &bytes.Buffer{}
	zr := lz4.NewReader(r)
	_, err := io.Copy(w, zr)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
