package compression

import (
	"errors"
	"fmt"
	"os"

	"github.com/pierrec/lz4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
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

func (m *L4zMagicNum) MagicNumMetadata() map[string]string {
	mn := make(map[string]string)
	mn["magic"] = fmt.Sprintf("%d", m.L)
	return mn
}

func (m *L4zMagicNum) GetMagicNumKeyValue(mnMap map[string]string) {
	if v, ok := mnMap["magic"]; ok {
		m.L = string_utils.IntStringParser(v)
	} else {
		panic(errors.New("no magic num provided"))
	}
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
	mn := L4zMagicNum{}
	mn.GetMagicNumKeyValue(p.Metadata)
	o, err := decompress(b, mn)
	err = inMemFs.MakeFileOut(p, o)
	p.DirIn = p.DirOut
	p.FnIn = p.FnOut
	if err != nil {
		log.Err(err).Msg("Lz4DecompressInMemFsFile: ")
		return inMemFs, err
	}
	return inMemFs, err
}

type L4zMagicNum struct {
	L int
}

func compress(toCompress []byte) ([]byte, L4zMagicNum, error) {
	compressed := make([]byte, len(toCompress))
	l, err := lz4.CompressBlock(toCompress, compressed, nil)
	if err != nil {
		panic(err)
	}
	mn := L4zMagicNum{L: len(toCompress)}
	return compressed[:l], mn, nil
}

func decompress(in []byte, mn L4zMagicNum) ([]byte, error) {
	decompressed := make([]byte, mn.L)
	l, err := lz4.UncompressBlock(in, decompressed)
	if err != nil {
		panic(err)
	}
	return decompressed[:l], err
}
